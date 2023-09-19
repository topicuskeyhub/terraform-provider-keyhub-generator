package model

import (
	"fmt"
	"log"

	"golang.org/x/exp/maps"
)

type restSimpleType struct {
	property             *RestProperty
	openapiType          string
	openapiFormat        string
	rsSchemaTemplateBase map[string]any
}

func NewRestSimpleType(property *RestProperty, name string, format string, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restSimpleType{
		property:             property,
		openapiType:          name,
		openapiFormat:        format,
		rsSchemaTemplateBase: rsSchemaTemplateBase,
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

func (t *restSimpleType) TFValueType() string {
	switch t.openapiType {
	case "boolean":
		return "basetypes.BoolValue"
	case "string":
		return "basetypes.StringValue"
	case "integer":
		return "basetypes.Int64Value"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) Complex() bool {
	return false
}

func (t *restSimpleType) NestedType() RestType {
	return nil
}

func (t *restSimpleType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restSimpleType) ToTKHAttrWithDiag() bool {
	return t.openapiType == "string" &&
		(t.openapiFormat == "date-time" || t.openapiFormat == "uuid" || t.openapiFormat == "date")
}

func (t *restSimpleType) ToTKHCustomCode() string {
	return ""
}

func (t *restSimpleType) TFAttrNeeded() bool {
	return false
}

func (t *restSimpleType) TKHToTF(value string, listItem bool) string {
	if listItem {
		switch t.openapiType {
		case "boolean":
			return "types.BoolValue(" + value + ")"
		case "string":
			if t.openapiFormat == "date-time" {
				return "timeToTF(" + value + ")"
			} else if t.openapiFormat == "uuid" || t.openapiFormat == "date" {
				return "types.StringValue(" + value + ".String())"
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
				return "timePointerToTF(" + value + ")"
			} else if t.openapiFormat == "uuid" || t.openapiFormat == "date" {
				return "stringerToTF(" + value + ")"
			}
			return "types.StringPointerValue(" + value + ")"
		case "integer":
			if t.openapiFormat == "int32" {
				return "types.Int64PointerValue(int32PToInt64P(" + value + "))"
			}
			return "types.Int64PointerValue(" + value + ")"
		default:
			log.Fatalf("Unknown simple type: %s", t.openapiType)
			return "error"
		}
	}
}

func (t *restSimpleType) TFToTKH(value string, listItem bool) string {
	if listItem {
		switch t.openapiType {
		case "boolean":
			return value + ".(basetypes.BoolValue).ValueBool()"
		case "string":
			if t.openapiFormat == "date-time" {
				return "tfToTime(" + value + ".(basetypes.StringValue))"
			} else if t.openapiFormat == "uuid" {
				return "parse(" + value + ".(basetypes.StringValue), uuid.Parse)"
			} else if t.openapiFormat == "date" {
				return "parse(" + value + ".(basetypes.StringValue), serialization.ParseDateOnly)"
			}
			return value + ".(basetypes.StringValue).ValueString()"
		case "integer":
			if t.openapiFormat == "int32" {
				return "int32(" + value + ".(basetypes.Int64Value).ValueInt64())"
			}
			return value + ".(basetypes.Int64Value).ValueInt64()"
		default:
			log.Fatalf("Unknown simple type: %s", t.openapiType)
			return "error"
		}
	} else {
		switch t.openapiType {
		case "boolean":
			return value + ".(basetypes.BoolValue).ValueBoolPointer()"
		case "string":
			if t.openapiFormat == "date-time" {
				return "tfToTimePointer(" + value + ".(basetypes.StringValue))"
			} else if t.openapiFormat == "uuid" {
				return "parsePointer(" + value + ".(basetypes.StringValue), uuid.Parse)"
			} else if t.openapiFormat == "date" {
				return "parsePointer2(" + value + ".(basetypes.StringValue), serialization.ParseDateOnly)"
			}
			return value + ".(basetypes.StringValue).ValueStringPointer()"
		case "integer":
			if t.openapiFormat == "int32" {
				return "int64PToInt32P(" + value + ".(basetypes.Int64Value).ValueInt64Pointer())"
			}
			return value + ".(basetypes.Int64Value).ValueInt64Pointer()"
		default:
			log.Fatalf("Unknown simple type: %s", t.openapiType)
			return "error"
		}
	}
}

func (t *restSimpleType) TKHToTFGuard() string {
	return ""
}

func (t *restSimpleType) TFToTKHGuard() string {
	return ""
}

func (t *restSimpleType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restSimpleType) SDKTypeName(listItem bool) string {
	var ret string
	switch t.openapiType {
	case "boolean":
		ret = "bool"
	case "string":
		if t.openapiFormat == "date-time" {
			ret = "time.Time"
		} else if t.openapiFormat == "uuid" {
			ret = "uuid.UUID"
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
	if !listItem {
		ret = "*" + ret
	}
	return ret
}

func (t *restSimpleType) SDKTypeConstructor() string {
	return "ERROR"
}

func (t *restSimpleType) DSSchemaTemplate() string {
	return "data_source_schema_attr_simple.go.tmpl"
}

func (t *restSimpleType) DSSchemaTemplateData() map[string]any {
	var attrType string
	switch t.openapiType {
	case "boolean":
		attrType = "dsschema.BoolAttribute"
	case "string":
		attrType = "dsschema.StringAttribute"
	case "integer":
		attrType = "dsschema.Int64Attribute"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		attrType = "error"
	}

	return map[string]any{
		"Type":     attrType,
		"Required": t.property.Name == "uuid",
	}
}

func (t *restSimpleType) RSSchemaTemplate() string {
	return "resource_schema_attr_simple.go.tmpl"
}

func (t *restSimpleType) RSSchemaTemplateData() map[string]any {
	var attrType string
	var planModifierType string
	var planModifierPkg string
	var defaultVal string
	switch t.openapiType {
	case "boolean":
		attrType = "rsschema.BoolAttribute"
		planModifierType = "planmodifier.Bool"
		planModifierPkg = "boolplanmodifier"
		defaultVal = fmt.Sprintf("booldefault.StaticBool(%v)", t.rsSchemaTemplateBase["Default"])
	case "string":
		attrType = "rsschema.StringAttribute"
		planModifierType = "planmodifier.String"
		planModifierPkg = "stringplanmodifier"
		defaultVal = fmt.Sprintf("stringdefault.StaticString(\"%v\")", t.rsSchemaTemplateBase["Default"])
	case "integer":
		attrType = "rsschema.Int64Attribute"
		planModifierType = "planmodifier.Int64"
		planModifierPkg = "int64planmodifier"
		defaultVal = fmt.Sprintf("int64default.StaticInt64(%v)", t.rsSchemaTemplateBase["Default"])
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		attrType = "error"
	}

	ret := map[string]any{
		"Type":             attrType,
		"PlanModifierType": planModifierType,
		"PlanModifierPkg":  planModifierPkg,
		"DefaultVal":       defaultVal,
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restSimpleType) DS() RestPropertyType {
	return NewRestSimpleType(t.property.DS(), t.openapiType, t.openapiFormat, t.rsSchemaTemplateBase)
}
