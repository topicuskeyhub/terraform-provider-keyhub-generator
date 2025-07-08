// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

type restFindByUUIDNestedItemType struct {
	nestedType         RestPropertyType
	parentPropertyName string
}

func NewFindByUUIDNestedItemType(nestedType RestPropertyType, parentPropertyName string) RestPropertyType {
	return &restFindByUUIDNestedItemType{
		nestedType:         nestedType,
		parentPropertyName: parentPropertyName,
	}
}

func (t *restFindByUUIDNestedItemType) MarkReachable() {
	t.nestedType.MarkReachable()
}

func (t *restFindByUUIDNestedItemType) PropertyNameSuffix() string {
	return "Uuid"
}

func (t *restFindByUUIDNestedItemType) FlattenMode() string {
	return "None"
}

func (t *restFindByUUIDNestedItemType) OrderMode() string {
	return "Object"
}

func (t *restFindByUUIDNestedItemType) TFName() string {
	return "types.Object"
}

func (t *restFindByUUIDNestedItemType) TFAttrType(inAdditionalObjects bool) string {
	var recurseCutOff = "false"

	// nestedAttrType := "objectAttrsType" + t.nestedType.NestedType().GoTypeName() + "UUIDItemType(" + recurseCutOff + ")"
	// return "types.ObjectType{AttrTypes: " + nestedAttrType + "}"
	return "types.ObjectType{AttrTypes: objectAttrsTypeRSFindByUUIDNestedItemType(" + recurseCutOff + ")}"
}

func (t *restFindByUUIDNestedItemType) TFValueType() string {
	return "basetypes.ObjectValue"
}

func (t *restFindByUUIDNestedItemType) TFValidatorType() string {
	return ""
}

func (t *restFindByUUIDNestedItemType) TFValidators() []string {
	return nil
}

func (t *restFindByUUIDNestedItemType) Complex() bool {
	return true
}

func (t *restFindByUUIDNestedItemType) RequiresReplace() bool {
	return false
}

func (t *restFindByUUIDNestedItemType) NestedType() RestType {
	return t.nestedType.NestedType()
}

func (t *restFindByUUIDNestedItemType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restFindByUUIDNestedItemType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restFindByUUIDNestedItemType) ToTKHCustomCode(buildType RestType) string {
	return ""
}

func (t *restFindByUUIDNestedItemType) TFAttrNeeded() bool {
	return false
}

func (t *restFindByUUIDNestedItemType) TKHToTF(value string, listItem bool) string {
	return "tkhToTFNestedFindByUuidType(" + t.nestedType.TKHToTF(value, false) + ")"
}

func (t *restFindByUUIDNestedItemType) TFToTKH(planValue string, configValue string, listItem bool) string {
	var tfPlanVal = "toObjectValue(" + planValue + ").Attributes()[\"uuid\"]"
	var tfConfigVal = "toObjectValue(" + configValue + ").Attributes()[\"uuid\"]"

	return t.nestedType.TFToTKH(tfPlanVal, tfConfigVal, listItem)
}

func (t *restFindByUUIDNestedItemType) TKHToTFGuard() string {
	return ""
}

func (t *restFindByUUIDNestedItemType) TFToTKHGuard() string {
	return ""
}

func (t *restFindByUUIDNestedItemType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restFindByUUIDNestedItemType) SDKInterfaceTypeName(listItem bool) string {
	return t.nestedType.NestedType().SDKInterfaceTypeName()
}

func (t *restFindByUUIDNestedItemType) SDKTypeConstructor() string {
	return t.nestedType.NestedType().SDKTypeConstructor()
}

func (t *restFindByUUIDNestedItemType) DSSchemaTemplate() string {
	return "data_source_schema_attr_nestedobject.go.tmpl"
}

func (t *restFindByUUIDNestedItemType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Type": "dsschema.ObjectAttribute",
	}
}

func (t *restFindByUUIDNestedItemType) RSSchemaTemplate() string {
	return "resource_schema_attr_nestedobject.go.tmpl"
}

func (t *restFindByUUIDNestedItemType) RSSchemaTemplateData() map[string]any {
	ret := map[string]any{
		"Type":             "rsschema.ObjectAttribute",
		"PlanModifierType": "planmodifier.Object",
		"PlanModifierPkg":  "objectplanmodifier",
		// "DefaultVal":       fmt.Sprintf("stringdefault.StaticString(\"%v\")", t.rsSchemaTemplateBase["Default"]),
	}
	// maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restFindByUUIDNestedItemType) DS() RestPropertyType {
	return t.nestedType.DS()
}
