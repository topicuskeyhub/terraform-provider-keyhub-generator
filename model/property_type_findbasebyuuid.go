// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"fmt"
)

type restFindBaseByUUIDObjectType struct {
	baseType *restFindByUUIDClassType
}

func NewFindBaseByUUIDObjectType(baseType *restFindByUUIDClassType) RestPropertyType {
	return &restFindBaseByUUIDObjectType{
		baseType: baseType,
	}
}

func (t *restFindBaseByUUIDObjectType) MarkReachable() {
	t.baseType.MarkReachable()
}

func (t *restFindBaseByUUIDObjectType) PropertyNameSuffix() string {
	return ""
}

func (t *restFindBaseByUUIDObjectType) FlattenMode() string {
	return "None"
}

func (t *restFindBaseByUUIDObjectType) OrderMode() string {
	return "None"
}

func (t *restFindBaseByUUIDObjectType) TFName() string {
	return "types.String"
}

func (t *restFindBaseByUUIDObjectType) TFAttrType(inAdditionalObjects bool) string {
	return "types.StringType"
}

func (t *restFindBaseByUUIDObjectType) TFValueType() string {
	return "basetypes.StringValue"
}

func (t *restFindBaseByUUIDObjectType) TFValidatorType() string {
	return "validator.String"
}

func (t *restFindBaseByUUIDObjectType) TFValidators() []string {
	return []string{`stringvalidator.RegexMatches(regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"), "The value must be a valid UUID"),`}
}

func (t *restFindBaseByUUIDObjectType) Complex() bool {
	return false
}

func (t *restFindBaseByUUIDObjectType) RequiresReplace() bool {
	return false
}

func (t *restFindBaseByUUIDObjectType) NestedType() RestType {
	return nil
}

func (t *restFindBaseByUUIDObjectType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restFindBaseByUUIDObjectType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restFindBaseByUUIDObjectType) ToTKHCustomCode(buildType RestType) string {
	typename := buildType.GoTypeName()
	superTypename := t.baseType.superClass.GoTypeName()
	return fmt.Sprintf("if val != nil {\n"+
		"dtype := tkh.GetTypeEscaped()\n"+
		"(*tkh.(*keyhubmodel.%s)).%s = *(val.(*keyhubmodel.%s))\n"+
		"tkh.SetTypeEscaped(dtype)\n}", typename, superTypename, superTypename)
}

func (t *restFindBaseByUUIDObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restFindBaseByUUIDObjectType) TKHToTF(value string, listItem bool) string {
	return "withUuidToTF(tkh)"
}

func (t *restFindBaseByUUIDObjectType) TFToTKH(value string, listItem bool) string {
	return "find" + t.baseType.superClass.GoTypeName() + "ByUUID(ctx, " + value + ".(basetypes.StringValue).ValueStringPointer())"
}

func (t *restFindBaseByUUIDObjectType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restFindBaseByUUIDObjectType) TKHToTFGuard() string {
	return ""
}

func (t *restFindBaseByUUIDObjectType) TFToTKHGuard() string {
	return ""
}

func (t *restFindBaseByUUIDObjectType) SDKTypeName(listItem bool) string {
	return "NONE"
}

func (t *restFindBaseByUUIDObjectType) SDKTypeConstructor() string {
	return "NONE"
}

func (t *restFindBaseByUUIDObjectType) DSSchemaTemplate() string {
	return "NONE"
}

func (t *restFindBaseByUUIDObjectType) DSSchemaTemplateData() map[string]any {
	return map[string]any{}
}

func (t *restFindBaseByUUIDObjectType) RSSchemaTemplate() string {
	return "resource_schema_attr_simple.go.tmpl"
}

func (t *restFindBaseByUUIDObjectType) RSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Type":             "rsschema.StringAttribute",
		"PlanModifierType": "planmodifier.String",
		"PlanModifierPkg":  "stringplanmodifier",
		"Mode":             "Required",
	}
}

func (t *restFindBaseByUUIDObjectType) DS() RestPropertyType {
	return nil
}
