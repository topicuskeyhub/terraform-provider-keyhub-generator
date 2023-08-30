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
			Name: name,
		}
		return ret
	} else {
		ret := &restClassType{
			SuperClass: superType,
			Name:       name,
		}
		types[name] = ret
		ret.Properties = buildProperties(ret, name, ownType, types)
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
		restProperty.Type = buildType(baseTypeName, name, property, types, restProperty)
		ret = append(ret, restProperty)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})
	return ret
}

func buildType(baseTypeName string, propertyName string, ref *openapi3.SchemaRef, types map[string]RestType, restProperty *RestProperty) RestPropertyType {
	schema := ref.Value
	if len(schema.AllOf) > 0 {
		if ref.Ref == "" {
			ref = schema.AllOf[0]
		}
		schema = schema.AllOf[0].Value
	}
	if schema.Type == "array" {
		return NewRestArrayType(buildType(baseTypeName, propertyName, schema.Items, types, restProperty))
	}
	if ref.Ref != "" && schema.Type == "string" && len(schema.Enum) > 0 {
		enum := getOrBuildTypeModel(types, refToName(ref.Ref), ref)
		return NewEnumPropertyType(enum)
	}
	if schema.Type == "boolean" || schema.Type == "integer" || schema.Type == "string" {
		return NewRestSimpleType(schema.Type, schema.Format)
	}
	if ref.Ref != "" && schema.Type == "object" && hasUUID(ref) {
		nested := getOrBuildTypeModel(types, refToName(ref.Ref), ref)
		if nested.Extends("Linkable") {
			return NewFindByUUIDObjectType(nested)
		}
	}
	if schema.Type == "object" || (len(ref.Value.AllOf) > 1 && ref.Value.AllOf[1].Value.Type == "object") {
		nestedTypeName := refToName(ref.Ref)
		if nestedTypeName == "" {
			nestedTypeName = baseTypeName + "_" + propertyName
		}
		nested := getOrBuildTypeModel(types, nestedTypeName, ref)
		return NewNestedObjectType(restProperty, nested)
	}

	log.Fatalf("Cannot construct a type for (%v+)", ref)
	return nil
}

func hasUUID(ref *openapi3.SchemaRef) bool {
	if _, ok := ref.Value.Properties["uuid"]; ok {
		return true
	}
	for _, part := range ref.Value.AllOf {
		if part.Ref == "" && hasUUID(part) {
			return true
		}
	}
	return false
}

func refToName(ref string) string {
	return ref[strings.LastIndex(ref, "/")+1:]
}
