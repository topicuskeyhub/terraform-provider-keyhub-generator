package model

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type restEnumPropertyType struct {
	enumType             *restEnumType
	rsSchemaTemplateBase map[string]any
}

func NewEnumPropertyType(enumType RestType, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restEnumPropertyType{
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

func (t *restEnumPropertyType) NestedType() RestType {
	return t.enumType
}

func (t *restEnumPropertyType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restEnumPropertyType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restEnumPropertyType) ToTKHCustomCode() string {
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

func (t *restEnumPropertyType) TFToTKH(value string, listItem bool) string {
	caster := "func(val any) " + t.SDKTypeName(listItem) + " { return *val.(*" + t.SDKTypeName(listItem) + ") }"
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

func (t *restEnumPropertyType) SDKTypeName(listItem bool) string {
	return t.enumType.SDKTypeName()
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
	return NewEnumPropertyType(t.enumType.DS(), t.rsSchemaTemplateBase)
}
