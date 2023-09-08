package model

import (
	"log"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

func BuildModel(openapi *openapi3.T) map[string]RestType {
	ret := make(map[string]RestType, 100)
	for name, schema := range openapi.Components.Schemas {
		if name == "RequestRange" {
			continue
		}
		getOrBuildTypeModel(ret, name, schema)
	}
	return ret
}

func getOrBuildTypeModel(types map[string]RestType, name string, schema *openapi3.SchemaRef) RestType {
	if ret, ok := types[name]; ok {
		return ret
	}

	var superType RestType
	if schema != nil && schema.Value.AllOf != nil {
		superType = getOrBuildTypeModel(types, refToName(schema.Value.AllOf[0].Ref), schema.Value.AllOf[0])
	}

	ownType := findOwnTypeSchema(schema)
	if ownType.Value.Type == "string" && len(ownType.Value.Enum) > 0 {
		ret := &restEnumType{
			suffix: "RS",
			name:   name,
		}
		return ret
	} else {
		ret := &restClassType{
			suffix:     "RS",
			superClass: superType,
			name:       name,
		}
		types[name] = ret
		ret.properties = buildProperties(ret, name, ownType, types)
		return ret
	}
}

func findOwnTypeSchema(schema *openapi3.SchemaRef) *openapi3.SchemaRef {
	if schema.Value.AllOf != nil {
		return schema.Value.AllOf[1]
	}
	return schema
}

func buildProperties(parent RestType, baseTypeName string, schema *openapi3.SchemaRef, types map[string]RestType) []*RestProperty {
	required := schema.Value.Required
	ret := make([]*RestProperty, 0)
	for name, property := range schema.Value.Properties {
		if name == "$type" {
			continue
		}
		rsSchemaTemplateBase := buildRSSchemaTemplateBase(schema, name)
		if name == "type" {
			if baseTypeName == "RestLink" || baseTypeName == "authPermission" {
				name = "typeEscaped"
			} else {
				name = baseTypeName + "Type"
			}
		} else if name == "vendor" {
			name = "vendorEscaped"
		}

		restProperty := &RestProperty{
			Parent:   parent,
			Name:     name,
			Required: slices.Contains(required, name),
		}
		restProperty.Type = buildType(baseTypeName, name, property, types, restProperty, rsSchemaTemplateBase)
		ret = append(ret, restProperty)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})
	return ret
}

func buildType(baseTypeName string, propertyName string, ref *openapi3.SchemaRef, types map[string]RestType, restProperty *RestProperty, rsSchemaTemplateBase map[string]any) RestPropertyType {
	schema := ref.Value
	if len(schema.AllOf) > 0 {
		if ref.Ref == "" {
			ref = schema.AllOf[0]
		}
		schema = schema.AllOf[0].Value
	}
	if schema.Type == "array" {
		return NewRestArrayType(buildType(baseTypeName, propertyName, schema.Items, types, restProperty, rsSchemaTemplateBase), rsSchemaTemplateBase)
	}
	if ref.Ref != "" && schema.Type == "string" && len(schema.Enum) > 0 {
		enum := getOrBuildTypeModel(types, refToName(ref.Ref), ref)
		return NewEnumPropertyType(enum, rsSchemaTemplateBase)
	}
	if schema.Type == "boolean" || schema.Type == "integer" || schema.Type == "string" {
		return NewRestSimpleType(restProperty, schema.Type, schema.Format, rsSchemaTemplateBase)
	}
	if isObject(ref) {
		nestedTypeName := refToName(ref.Ref)
		if nestedTypeName == "" {
			nestedTypeName = baseTypeName + "_" + propertyName
		}
		nested := getOrBuildTypeModel(types, nestedTypeName, ref)
		ret := NewNestedObjectType(restProperty, nested, rsSchemaTemplateBase)
		if ref.Ref != "" && hasUUID(ref) && nested.Extends("Linkable") {
			if strings.HasSuffix(nestedTypeName, "Primer") {
				ret = NewFindByUUIDObjectType(ret, rsSchemaTemplateBase)
			}
		}
		return ret
	}

	log.Fatalf("Cannot construct a type for (%v+)", ref)
	return nil
}

func hasUUID(ref *openapi3.SchemaRef) bool {
	if _, ok := ref.Value.Properties["uuid"]; ok {
		return true
	}
	for _, part := range ref.Value.AllOf {
		if hasUUID(part) {
			return true
		}
	}
	return false
}

func isObject(ref *openapi3.SchemaRef) bool {
	if ref.Value.Type == "object" {
		return true
	}
	for _, part := range ref.Value.AllOf {
		if isObject(part) {
			return true
		}
	}
	return false
}

func refToName(ref string) string {
	return ref[strings.LastIndex(ref, "/")+1:]
}

/*
Computed:
Read only field, returned from backend

Computed_UseStateForUnknown
Immutable field, returned from backend

Computed_Optional (Not support for now)
Optional field with a default value determined by the backend if not provided

Optional
Optional field that can be left empty

Optional_Default
Optional or required field that has a default value

Required
Required field that does not have a default value
*/
func buildRSSchemaTemplateBase(ref *openapi3.SchemaRef, propertyName string) map[string]any {
	required := slices.Contains(ref.Value.Required, propertyName)
	property := ref.Value.Properties[propertyName]
	if property.Value.AllOf != nil && len(property.Value.AllOf) > 0 {
		property = property.Value.AllOf[1]
	}
	defaultVal := property.Value.Extensions["x-tkh-default"]
	readOnly := property.Value.ReadOnly
	immutable := property.Value.Extensions["x-tkh-immutable"] != nil && property.Value.Extensions["x-tkh-immutable"].(bool)
	createOnly := property.Value.Extensions["x-tkh-create-only"] != nil && property.Value.Extensions["x-tkh-create-only"].(bool)
	if immutable {
		return map[string]any{
			"Mode": "Computed_UseStateForUnknown",
		}
	}
	if readOnly && !createOnly {
		return map[string]any{
			"Mode": "Computed",
		}
	}
	if defaultVal != nil {
		return map[string]any{
			"Mode":    "Optional_Default",
			"Default": defaultVal,
		}
	}
	if required {
		return map[string]any{
			"Mode": "Required",
		}
	}
	return map[string]any{
		"Mode": "Optional",
	}
}
