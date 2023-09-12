package model

type restFindBaseByUUIDObjectType struct {
	baseType RestType
}

func NewFindBaseByUUIDObjectType(baseType RestType) RestPropertyType {
	return &restFindBaseByUUIDObjectType{
		baseType: baseType,
	}
}

func (t *restFindBaseByUUIDObjectType) PropertyNameSuffix() string {
	return "Uuid"
}

func (t *restFindBaseByUUIDObjectType) TFName() string {
	return "types.String"
}

func (t *restFindBaseByUUIDObjectType) TFAttrType() string {
	return "types.StringType"
}

func (t *restFindBaseByUUIDObjectType) TFValueType() string {
	return "basetypes.StringValue"
}

func (t *restFindBaseByUUIDObjectType) Complex() bool {
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

func (t *restFindBaseByUUIDObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restFindBaseByUUIDObjectType) TKHToTF(value string, listItem bool) string {
	return "withUuidToTF(tkh)"
}

func (t *restFindBaseByUUIDObjectType) TFToTKH(value string, listItem bool) string {
	return "find" + t.baseType.GoTypeName() + "ByUUID(ctx, " + value + ".(basetypes.StringValue).ValueStringPointer())"
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
