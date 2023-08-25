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

func (t *restSimpleType) TKHToTF(value string, list bool) string {
	if list {
		switch t.openapiType {
		case "boolean":
			return "types.BoolValue(" + value + ")"
		case "string":
			if t.openapiFormat == "date-time" {
				return "TimeToTF(" + value + ")"
			}
			return "types.StringValue(" + value + ")"
		case "integer":
			if t.openapiFormat == "int32" {
				return "types.Int64Value(int64(" + value + "))"
			}
			return "types.Int64Value(" + value + ")"
		default:
			log.Fatalf("Unknown simple type: %s", t.openapiType)
			return "error"
		}
	} else {
		switch t.openapiType {
		case "boolean":
			return "types.BoolPointerValue(" + value + ")"
		case "string":
			if t.openapiFormat == "date-time" {
				return "TimePointerToTF(" + value + ")"
			}
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
}

func (t *restSimpleType) SDKTypeName(list bool) string {
	var ret string
	switch t.openapiType {
	case "boolean":
		ret = "bool"
	case "string":
		if t.openapiFormat == "date-time" {
			ret = "time.Time"
		} else {
			ret = "string"
		}
	case "integer":
		if t.openapiFormat == "int32" {
			ret = "int32"
		} else {
			ret = "int64"
		}
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
	if !list {
		ret = "*" + ret
	}
	return ret
}
