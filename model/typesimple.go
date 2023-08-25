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

func (t *restSimpleType) TFAttrType() string {
	switch t.openapiType {
	case "boolean":
		return "types.BoolType"
	case "string":
		return "types.StringType"
	case "integer":
		return "types.Int64Type"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) NestedType() RestType {
	return nil
}

func (t *restSimpleType) TFAttrWithDiag() bool {
	return false
}

func (t *restSimpleType) TFAttrNeeded() bool {
	return false
}

func (t *restSimpleType) TKHToTF(value string) string {
	switch t.openapiType {
	case "boolean":
		return "types.BoolPointerValue(" + value + ")"
	case "string":
		return "types.StringPointerValue(" + value + ")"
	case "integer":
		if t.openapiFormat == "int32" {
			return "types.Int64PointerValue(Int32PToInt64P(" + value + "))"
		}
		return "types.Int64PointerValue(" + value + ")"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) SDKTypeName() string {
	switch t.openapiType {
	case "boolean":
		return "*bool"
	case "string":
		return "*string"
	case "integer":
		return "*int64"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}
