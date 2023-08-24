package model

import (
	"log"
)

type restSimpleType struct {
	openapiType   string
	openapiFormat string
}

func NewRestSimpleType(name string, format string) RestPropertyType {
	return &restSimpleType{
		openapiType:   name,
		openapiFormat: format,
	}
}

func (t *restSimpleType) PropertyNameSuffix() string {
	return ""
}

func (t *restSimpleType) TFName() string {
	switch t.openapiType {
	case "boolean":
		return "types.Bool"
	case "string":
		return "types.String"
	case "integer":
		return "types.Int64"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) NestedType() *RestType {
	return nil
}
