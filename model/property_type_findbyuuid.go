// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"fmt"

	"golang.org/x/exp/maps"
)

type restFindByUUIDObjectType struct {
	nestedType           RestPropertyType
	rsSchemaTemplateBase map[string]any
}

func NewFindByUUIDObjectType(nestedType RestPropertyType, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restFindByUUIDObjectType{
		nestedType:           nestedType,
		rsSchemaTemplateBase: rsSchemaTemplateBase,
	}
}

func (t *restFindByUUIDObjectType) MarkReachable() {
	t.nestedType.MarkReachable()
}

func (t *restFindByUUIDObjectType) PropertyNameSuffix() string {
	return "Uuid"
}

func (t *restFindByUUIDObjectType) FlattenMode() string {
	return "None"
}

func (t *restFindByUUIDObjectType) OrderMode() string {
	return "None"
}

func (t *restFindByUUIDObjectType) TFName() string {
	return "types.String"
}

func (t *restFindByUUIDObjectType) TFAttrType(inAdditionalObjects bool) string {
	return "types.StringType"
}

func (t *restFindByUUIDObjectType) TFValueType() string {
	return "basetypes.StringValue"
}

func (t *restFindByUUIDObjectType) TFValidatorType() string {
	return "validator.String"
}

func (t *restFindByUUIDObjectType) TFValidators() []string {
	return []string{`stringvalidator.RegexMatches(regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"), "The value must be a valid UUID"),`}
}

func (t *restFindByUUIDObjectType) Complex() bool {
	return false
}

func (t *restFindByUUIDObjectType) RequiresReplace() bool {
	return false
}

func (t *restFindByUUIDObjectType) NestedType() RestType {
	return t.nestedType.NestedType()
}

func (t *restFindByUUIDObjectType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restFindByUUIDObjectType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restFindByUUIDObjectType) ToTKHCustomCode(buildType RestType) string {
	return ""
}

func (t *restFindByUUIDObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restFindByUUIDObjectType) TKHToTF(value string, listItem bool) string {
	return "withUuidToTF(" + value + ")"
}

func (t *restFindByUUIDObjectType) TFToTKH(value string, listItem bool) string {
	return "find" + t.nestedType.NestedType().GoTypeName() + "ByUUID(ctx, " + value + ".(basetypes.StringValue).ValueStringPointer())"
}

func (t *restFindByUUIDObjectType) TKHToTFGuard() string {
	return ""
}

func (t *restFindByUUIDObjectType) TFToTKHGuard() string {
	return ""
}

func (t *restFindByUUIDObjectType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restFindByUUIDObjectType) SDKTypeName(listItem bool) string {
	return t.nestedType.NestedType().SDKTypeName()
}

func (t *restFindByUUIDObjectType) SDKTypeConstructor() string {
	return t.nestedType.NestedType().SDKTypeConstructor()
}

func (t *restFindByUUIDObjectType) DSSchemaTemplate() string {
	return "data_source_schema_attr_simple.go.tmpl"
}

func (t *restFindByUUIDObjectType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Type": "dsschema.StringAttribute",
	}
}

func (t *restFindByUUIDObjectType) RSSchemaTemplate() string {
	return "resource_schema_attr_simple.go.tmpl"
}

func (t *restFindByUUIDObjectType) RSSchemaTemplateData() map[string]any {
	ret := map[string]any{
		"Type":             "rsschema.StringAttribute",
		"PlanModifierType": "planmodifier.String",
		"PlanModifierPkg":  "stringplanmodifier",
		"DefaultVal":       fmt.Sprintf("stringdefault.StaticString(\"%v\")", t.rsSchemaTemplateBase["Default"]),
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restFindByUUIDObjectType) DS() RestPropertyType {
	return t.nestedType.DS()
}
