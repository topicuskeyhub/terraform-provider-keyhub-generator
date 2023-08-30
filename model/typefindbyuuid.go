package model

type restFindByUUIDObjectType struct {
	nestedType RestType
}

func NewFindByUUIDObjectType(nestedType RestType) RestPropertyType {
	return &restFindByUUIDObjectType{
		nestedType: nestedType,
	}
}

func (t *restFindByUUIDObjectType) PropertyNameSuffix() string {
	return "Uuid"
}

func (t *restFindByUUIDObjectType) TFName() string {
	return "types.String"
}

func (t *restFindByUUIDObjectType) TFAttrType() string {
	return "types.StringType"
}

func (t *restFindByUUIDObjectType) TFValueType() string {
	return "basetypes.StringValue"
}

func (t *restFindByUUIDObjectType) Complex() bool {
	return false
}

func (t *restFindByUUIDObjectType) NestedType() RestType {
	return nil
}

func (t *restFindByUUIDObjectType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restFindByUUIDObjectType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restFindByUUIDObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restFindByUUIDObjectType) TKHToTF(value string, listItem bool) string {
	return "withUuidToTF(" + value + ")"
}

func (t *restFindByUUIDObjectType) TFToTKH(value string, listItem bool) string {
	return "find" + t.nestedType.GoTypeName() + "ByUUID(ctx, " + value + ".(basetypes.StringValue).ValueStringPointer())"
}

func (t *restFindByUUIDObjectType) SDKTypeName(listItem bool) string {
	return t.nestedType.SDKTypeName()
}

func (t *restFindByUUIDObjectType) SDKTypeConstructor() string {
	return t.nestedType.SDKTypeConstructor()
}

func (t *restFindByUUIDObjectType) DSSchemaTemplate() string {
	return "data_source_schema_attr_simple.go.tmpl"
}

func (t *restFindByUUIDObjectType) DSSchemaTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Type": "dsschema.StringAttribute",
	}
}
