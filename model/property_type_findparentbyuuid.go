// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

type restFindParentByUUIDObjectType struct {
}

func NewFindParentByUUIDObjectType() RestPropertyType {
	return &restFindParentByUUIDObjectType{}
}

func (t *restFindParentByUUIDObjectType) MarkReachable() {
}

func (t *restFindParentByUUIDObjectType) PropertyNameSuffix() string {
	return ""
}

func (t *restFindParentByUUIDObjectType) FlattenMode() string {
	return "None"
}

func (t *restFindParentByUUIDObjectType) OrderMode() string {
	return "None"
}

func (t *restFindParentByUUIDObjectType) TFName() string {
	return "types.String"
}

func (t *restFindParentByUUIDObjectType) TFAttrType(inAdditionalObjects bool) string {
	return "types.StringType"
}

func (t *restFindParentByUUIDObjectType) TFValueType() string {
	return "basetypes.StringValue"
}

func (t *restFindParentByUUIDObjectType) TFValidatorType() string {
	return "validator.String"
}

func (t *restFindParentByUUIDObjectType) TFValidators() []string {
	return []string{`stringvalidator.RegexMatches(regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"), "The value must be a valid UUID"),`}
}

func (t *restFindParentByUUIDObjectType) Complex() bool {
	return false
}

func (t *restFindParentByUUIDObjectType) RequiresReplace() bool {
	return true
}

func (t *restFindParentByUUIDObjectType) NestedType() RestType {
	return nil
}

func (t *restFindParentByUUIDObjectType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restFindParentByUUIDObjectType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restFindParentByUUIDObjectType) ToTKHCustomCode(buildType RestType) string {
	return ""
}

func (t *restFindParentByUUIDObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restFindParentByUUIDObjectType) TKHToTF(value string, listItem bool) string {
	return "types.StringNull()"
}

func (t *restFindParentByUUIDObjectType) TFToTKH(planValue string, configValue string, listItem bool) string {
	return ""
}

func (t *restFindParentByUUIDObjectType) TKHGetter(propertyName string) string {
	return ""
}

func (t *restFindParentByUUIDObjectType) TKHToTFGuard() string {
	return ""
}

func (t *restFindParentByUUIDObjectType) TFToTKHGuard() string {
	return ""
}

func (t *restFindParentByUUIDObjectType) SDKTypeName(listItem bool) string {
	return "NONE"
}

func (t *restFindParentByUUIDObjectType) SDKTypeConstructor() string {
	return "NONE"
}

func (t *restFindParentByUUIDObjectType) DSSchemaTemplate() string {
	return "data_source_schema_attr_simple.go.tmpl"
}

func (t *restFindParentByUUIDObjectType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Type":     "dsschema.StringAttribute",
		"Required": true,
	}
}

func (t *restFindParentByUUIDObjectType) RSSchemaTemplate() string {
	return "resource_schema_attr_simple.go.tmpl"
}

func (t *restFindParentByUUIDObjectType) RSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Type":             "rsschema.StringAttribute",
		"PlanModifierType": "planmodifier.String",
		"PlanModifierPkg":  "stringplanmodifier",
		"Mode":             "Required",
	}
}

func (t *restFindParentByUUIDObjectType) DS() RestPropertyType {
	return NewFindParentByUUIDObjectType()
}
