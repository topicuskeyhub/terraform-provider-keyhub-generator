package model

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

func BuildModel(openapi *openapi3.T) map[string]*RestType {
	ret := make(map[string]*RestType, 100)
	for name, schema := range openapi.Components.Schemas {
		ret[name] = buildTypeModel(name, schema)
	}
	for name, typeModel := range ret {
		schema := openapi.Components.Schemas[name]
		if schema.Value.AllOf != nil {
			ref := schema.Value.AllOf[0].Ref
			superClassName := ref[strings.LastIndex(ref, "/")+1:]
			typeModel.SuperClass = ret[superClassName]
		}
	}
	return ret
}

func buildTypeModel(name string, schema *openapi3.SchemaRef) *RestType {
	ownType := findOwnTypeSchema(schema)
	ret := &RestType{
		Name:       name,
		Properties: buildProperties(ownType),
	}
	return ret
}

func findOwnTypeSchema(schema *openapi3.SchemaRef) *openapi3.SchemaRef {
	if schema.Value.AllOf != nil {
		return schema.Value.AllOf[1]
	}
	return schema
}

func buildProperties(schema *openapi3.SchemaRef) []*RestProperty {
	required := schema.Value.Required
	ret := make([]*RestProperty, 0)
	for name, property := range schema.Value.Properties {
		restProperty := &RestProperty{
			Name:     name,
			Type:     property.Value.Type,
			Required: slices.Contains(required, name),
		}
		ret = append(ret, restProperty)
	}
	return ret
}
