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
		getOrBuildTypeModel(ret, name, schema)
	}
	for name, typeModel := range ret {
		schema := openapi.Components.Schemas[name]
		if schema != nil && schema.Value.AllOf != nil {
			superClassName := refToName(schema.Value.AllOf[0].Ref)
			typeModel.(*restClassType).SuperClass = ret[superClassName]
		}
	}
	return ret
}

func getOrBuildTypeModel(types map[string]RestType, name string, schema *openapi3.SchemaRef) RestType {
	if ret, ok := types[name]; ok {
		return ret
	}

	ownType := findOwnTypeSchema(schema)
	ret := &restClassType{
		Name: name,
	}
	types[name] = ret
	if name == "Linkable" {
		ret.Properties = make([]*RestProperty, 0)
	} else {
		ret.Properties = buildProperties(name, ownType, types)
	}
	return ret
}

func findOwnTypeSchema(schema *openapi3.SchemaRef) *openapi3.SchemaRef {
	if schema.Value.AllOf != nil {
		return schema.Value.AllOf[1]
	}
	return schema
}

func buildProperties(baseTypeName string, schema *openapi3.SchemaRef, types map[string]RestType) []*RestProperty {
	required := schema.Value.Required
	ret := make([]*RestProperty, 0)
	for name, property := range schema.Value.Properties {
		if name == "$type" {
			continue
		}

		restProperty := &RestProperty{
			Name:     name,
			Type:     buildType(baseTypeName, name, property, types),
			Required: slices.Contains(required, name),
		}
		ret = append(ret, restProperty)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})
	return ret
}

func buildType(baseTypeName string, propertyName string, ref *openapi3.SchemaRef, types map[string]RestType) RestPropertyType {
	schema := ref.Value
	if len(schema.AllOf) > 0 {
		schema = schema.AllOf[0].Value
	}
	if schema.Type == "array" {
		return NewRestArrayType(buildType(baseTypeName, propertyName, schema.Items, types))
	}
	if schema.Type == "boolean" || schema.Type == "integer" || schema.Type == "string" {
		return NewRestSimpleType(schema.Type, schema.Format)
	}
	if ref.Ref != "" && schema.Type == "object" && hasUUID(ref) {
		nested := getOrBuildTypeModel(types, refToName(ref.Ref), ref)
		return NewFindByUUIDObjectType(nested)
	}
	if schema.Type == "object" || (len(ref.Value.AllOf) > 1 && ref.Value.AllOf[1].Value.Type == "object") {
		nestedTypeName := refToName(ref.Ref)
		if nestedTypeName == "" {
			nestedTypeName = baseTypeName + "_" + propertyName
		}
		nested := getOrBuildTypeModel(types, nestedTypeName, ref)
		return NewNestedObjectType(nested)
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
