// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"fmt"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/maps"
)

type restSimpleType struct {
	property             *RestProperty
	openapiType          *openapi3.Types
	openapiSchema        *openapi3.Schema
	rsSchemaTemplateBase map[string]any
}

func NewRestSimpleType(property *RestProperty, schema *openapi3.Schema, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restSimpleType{
		property:             property,
		openapiType:          schema.Type,
		openapiSchema:        schema,
		rsSchemaTemplateBase: rsSchemaTemplateBase,
	}
}

func (t *restSimpleType) MarkReachable() {
}

func (t *restSimpleType) PropertyNameSuffix() string {
	return ""
}

func (t *restSimpleType) FlattenMode() string {
	return "None"
}

func (t *restSimpleType) OrderMode() string {
	return "None"
}

func (t *restSimpleType) TFName() string {
	switch {
	case t.openapiType.Is("boolean"):
		return "types.Bool"
	case t.openapiType.Is("string"):
		return "types.String"
	case t.openapiType.Is("integer"):
		return "types.Int64"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) TFAttrType(inAdditionalObjects bool) string {
	switch {
	case t.openapiType.Is("boolean"):
		return "types.BoolType"
	case t.openapiType.Is("string"):
		return "types.StringType"
	case t.openapiType.Is("integer"):
		return "types.Int64Type"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) TFValueType() string {
	switch {
	case t.openapiType.Is("boolean"):
		return "basetypes.BoolValue"
	case t.openapiType.Is("string"):
		return "basetypes.StringValue"
	case t.openapiType.Is("integer"):
		return "basetypes.Int64Value"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) TFValidatorType() string {
	switch {
	case t.openapiType.Is("boolean"):
		return "validator.Bool"
	case t.openapiType.Is("string"):
		return "validator.String"
	case t.openapiType.Is("integer"):
		return "validator.Int64"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		return "error"
	}
}

func (t *restSimpleType) TFValidators() []string {
	validators := make([]string, 0)
	if t.openapiType.Is("string") {
		minLength := t.openapiSchema.MinLength
		maxLength := t.openapiSchema.MaxLength
		if maxLength != nil {
			validators = append(validators, fmt.Sprintf("stringvalidator.UTF8LengthBetween(%d, %d),", minLength, *maxLength))
		}
		if t.property.Name == "uuid" {
			validators = append(validators, `stringvalidator.RegexMatches(regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"), "The value must be a valid UUID"),`)
		}
	} else if t.openapiType.Is("integer") {
		min := t.openapiSchema.Min
		max := t.openapiSchema.Max
		if min != nil {
			if max != nil {
				validators = append(validators, fmt.Sprintf("int64validator.Between(%d, %d),", int(*min), int(*max)))
			} else {
				validators = append(validators, fmt.Sprintf("int64validator.AtLeast(%d),", int(*min)))
			}
		} else if max != nil {
			validators = append(validators, fmt.Sprintf("int64validator.AtMost(%d),", int(*max)))
		}
	}
	return validators
}

func (t *restSimpleType) Complex() bool {
	return false
}

func (t *restSimpleType) RequiresReplace() bool {
	return false
}

func (t *restSimpleType) NestedType() RestType {
	return nil
}

func (t *restSimpleType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restSimpleType) ToTKHAttrWithDiag() bool {
	openapiFormat := t.openapiSchema.Format
	return t.openapiType.Is("string") &&
		(openapiFormat == "date-time" || openapiFormat == "uuid" || openapiFormat == "date")
}

func (t *restSimpleType) ToTKHCustomCode(buildType RestType) string {
	return ""
}

func (t *restSimpleType) TFAttrNeeded() bool {
	return false
}

func (t *restSimpleType) TKHToTF(value string, listItem bool) string {
	openapiFormat := t.openapiSchema.Format
	if listItem {
		switch {
		case t.openapiType.Is("boolean"):
			return "types.BoolValue(" + value + ")"
		case t.openapiType.Is("string"):
			if openapiFormat == "date-time" {
				return "timeToTF(" + value + ")"
			} else if openapiFormat == "uuid" || openapiFormat == "date" {
				return "types.StringValue(" + value + ".String())"
			} else if openapiFormat == "byte" {
				return "types.StringValue(string(" + value + "))"
			}
			return "types.StringValue(" + value + ")"
		case t.openapiType.Is("integer"):
			if openapiFormat == "int32" {
				return "types.Int64Value(int64(" + value + "))"
			}
			return "types.Int64Value(" + value + ")"
		default:
			log.Fatalf("Unknown simple type: %s", t.openapiType)
			return "error"
		}
	} else {
		switch {
		case t.openapiType.Is("boolean"):
			return "types.BoolPointerValue(" + value + ")"
		case t.openapiType.Is("string"):
			if openapiFormat == "date-time" {
				return "timePointerToTF(" + value + ")"
			} else if openapiFormat == "uuid" || openapiFormat == "date" {
				return "stringerToTF(" + value + ")"
			} else if openapiFormat == "byte" {
				return "types.StringValue(string(" + value + "))"
			}
			return "types.StringPointerValue(" + value + ")"
		case t.openapiType.Is("integer"):
			if openapiFormat == "int32" {
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
	openapiFormat := t.openapiSchema.Format
	if listItem {
		switch {
		case t.openapiType.Is("boolean"):
			return value + ".(basetypes.BoolValue).ValueBool()"
		case t.openapiType.Is("string"):
			if openapiFormat == "date-time" {
				return "tfToTime(" + value + ".(basetypes.StringValue))"
			} else if openapiFormat == "uuid" {
				return "parse(" + value + ".(basetypes.StringValue), uuid.Parse)"
			} else if openapiFormat == "date" {
				return "parse(" + value + ".(basetypes.StringValue), serialization.ParseDateOnly)"
			} else if openapiFormat == "byte" {
				return "[]byte(" + value + ".(basetypes.StringValue).ValueString())"
			}
			return value + ".(basetypes.StringValue).ValueString()"
		case t.openapiType.Is("integer"):
			if openapiFormat == "int32" {
				return "int32(" + value + ".(basetypes.Int64Value).ValueInt64())"
			}
			return value + ".(basetypes.Int64Value).ValueInt64()"
		default:
			log.Fatalf("Unknown simple type: %s", t.openapiType)
			return "error"
		}
	} else {
		switch {
		case t.openapiType.Is("boolean"):
			return "tfToBooleanPointer(" + value + ")"
		case t.openapiType.Is("string"):
			if openapiFormat == "date-time" {
				return "tfToTimePointer(" + value + ".(basetypes.StringValue))"
			} else if openapiFormat == "uuid" {
				return "parsePointer(" + value + ".(basetypes.StringValue), uuid.Parse)"
			} else if openapiFormat == "date" {
				return "parsePointer2(" + value + ".(basetypes.StringValue), serialization.ParseDateOnly)"
			} else if openapiFormat == "byte" {
				return "[]byte(" + value + ".(basetypes.StringValue).ValueString())"
			}
			return "tfToStringPointer(" + value + ")"
		case t.openapiType.Is("integer"):
			if openapiFormat == "int32" {
				return "int64PToInt32P(tfToInt64Pointer(" + value + "))"
			}
			return "tfToInt64Pointer(" + value + ")"
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
	openapiFormat := t.openapiSchema.Format
	var ret string
	switch {
	case t.openapiType.Is("boolean"):
		ret = "bool"
	case t.openapiType.Is("string"):
		if openapiFormat == "date-time" {
			ret = "time.Time"
		} else if openapiFormat == "uuid" {
			ret = "uuid.UUID"
		} else if openapiFormat == "byte" {
			ret = "[]byte"
		} else {
			ret = "string"
		}
	case t.openapiType.Is("integer"):
		if openapiFormat == "int32" {
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
	switch {
	case t.openapiType.Is("boolean"):
		attrType = "dsschema.BoolAttribute"
	case t.openapiType.Is("string"):
		attrType = "dsschema.StringAttribute"
	case t.openapiType.Is("integer"):
		attrType = "dsschema.Int64Attribute"
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		attrType = "error"
	}
	sensitive := t.openapiSchema.Extensions["x-tkh-sensitive"] != nil && t.openapiSchema.Extensions["x-tkh-sensitive"].(bool)

	return map[string]any{
		"Type":      attrType,
		"Required":  t.property.Name == "uuid",
		"Sensitive": sensitive,
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
	switch {
	case t.openapiType.Is("boolean"):
		attrType = "rsschema.BoolAttribute"
		planModifierType = "planmodifier.Bool"
		planModifierPkg = "boolplanmodifier"
		defaultVal = fmt.Sprintf("booldefault.StaticBool(%v)", t.rsSchemaTemplateBase["Default"])
	case t.openapiType.Is("string"):
		attrType = "rsschema.StringAttribute"
		planModifierType = "planmodifier.String"
		planModifierPkg = "stringplanmodifier"
		defaultVal = fmt.Sprintf("stringdefault.StaticString(\"%v\")", t.rsSchemaTemplateBase["Default"])
	case t.openapiType.Is("integer"):
		attrType = "rsschema.Int64Attribute"
		planModifierType = "planmodifier.Int64"
		planModifierPkg = "int64planmodifier"
		defaultVal = fmt.Sprintf("int64default.StaticInt64(%v)", t.rsSchemaTemplateBase["Default"])
	default:
		log.Fatalf("Unknown simple type: %s", t.openapiType)
		attrType = "error"
	}
	sensitive := t.openapiSchema.Extensions["x-tkh-sensitive"] != nil && t.openapiSchema.Extensions["x-tkh-sensitive"].(bool)

	ret := map[string]any{
		"Type":             attrType,
		"PlanModifierType": planModifierType,
		"PlanModifierPkg":  planModifierPkg,
		"DefaultVal":       defaultVal,
		"Sensitive":        sensitive,
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restSimpleType) DS() RestPropertyType {
	return NewRestSimpleType(t.property.DS(), t.openapiSchema, t.rsSchemaTemplateBase)
}
