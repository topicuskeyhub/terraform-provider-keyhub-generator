// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type restEnumPropertyType struct {
	property             *RestProperty
	enumType             *restEnumType
	rsSchemaTemplateBase map[string]any
}

func NewEnumPropertyType(property *RestProperty, enumType RestType, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restEnumPropertyType{
		property:             property,
		enumType:             enumType.(*restEnumType),
		rsSchemaTemplateBase: rsSchemaTemplateBase,
	}
}

func (t *restEnumPropertyType) MarkReachable() {
	t.enumType.MarkReachable()
}

func (t *restEnumPropertyType) PropertyNameSuffix() string {
	return ""
}

func (t *restEnumPropertyType) FlattenMode() string {
	return "None"
}

func (t *restEnumPropertyType) OrderMode() string {
	return "None"
}

func (t *restEnumPropertyType) TFName() string {
	return "types.String"
}

func (t *restEnumPropertyType) TFAttrType(inAdditionalObjects bool) string {
	return "types.StringType"
}

func (t *restEnumPropertyType) TFValueType() string {
	return "basetypes.StringValue"
}

func (t *restEnumPropertyType) TFValidatorType() string {
	return "validator.String"
}

func (t *restEnumPropertyType) TFValidators() []string {
	var sb strings.Builder
	sb.WriteString("stringvalidator.OneOf(\n")
	for _, name := range t.enumType.values {
		sb.WriteString(`"`)
		sb.WriteString(fmt.Sprint(name))
		sb.WriteString(`",`)
	}
	sb.WriteString("\n),")
	return []string{sb.String()}
}

func (t *restEnumPropertyType) Complex() bool {
	return false
}

func (t *restEnumPropertyType) RequiresReplace() bool {
	return false
}

func (t *restEnumPropertyType) NestedType() RestType {
	return t.enumType
}

func (t *restEnumPropertyType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restEnumPropertyType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restEnumPropertyType) ToTKHCustomCode(baseType RestType) string {
	return ""
}

func (t *restEnumPropertyType) TFAttrNeeded() bool {
	return false
}

func (t *restEnumPropertyType) TKHToTF(value string, listItem bool) string {
	if listItem {
		return "types.StringValue(" + value + ".String())"
	}
	return "stringerToTF(" + value + ")"
}

func (t *restEnumPropertyType) TFToTKH(planValue string, configValue string, listItem bool) string {
	var value string
	if t.property.IsValueFromConfig() {
		value = configValue
	} else {
		value = planValue
	}

	caster := "func(val any) " + t.SDKInterfaceTypeName(listItem) + " { return *val.(*" + t.SDKInterfaceTypeName(listItem) + ") }"
	if listItem {
		return "parseCast(" + value + ".(basetypes.StringValue), " + t.SDKTypeConstructor() + ", " + caster + ")"
	}
	return "parseCastPointer(" + value + ".(basetypes.StringValue), " + t.SDKTypeConstructor() + ", " + caster + ")"
}

func (t *restEnumPropertyType) TKHToTFGuard() string {
	return ""
}

func (t *restEnumPropertyType) TFToTKHGuard() string {
	return ""
}

func (t *restEnumPropertyType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restEnumPropertyType) SDKInterfaceTypeName(listItem bool) string {
	return t.enumType.SDKInterfaceTypeName()
}

func (t *restEnumPropertyType) SDKTypeConstructor() string {
	return t.enumType.SDKTypeConstructor()
}

func (t *restEnumPropertyType) DSSchemaTemplate() string {
	return "data_source_schema_attr_simple.go.tmpl"
}

func (t *restEnumPropertyType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Type": "dsschema.StringAttribute",
	}
}

func (t *restEnumPropertyType) RSSchemaTemplate() string {
	return "resource_schema_attr_simple.go.tmpl"
}

func (t *restEnumPropertyType) RSSchemaTemplateData() map[string]any {
	ret := map[string]any{
		"Type":             "rsschema.StringAttribute",
		"PlanModifierType": "planmodifier.String",
		"PlanModifierPkg":  "stringplanmodifier",
		"DefaultVal":       fmt.Sprintf("stringdefault.StaticString(\"%v\")", t.rsSchemaTemplateBase["Default"]),
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restEnumPropertyType) DS() RestPropertyType {
	return NewEnumPropertyType(t.property.DS(), t.enumType.DS(), t.rsSchemaTemplateBase)
}
