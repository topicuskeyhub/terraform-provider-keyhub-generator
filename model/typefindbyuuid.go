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

func (t *restFindByUUIDObjectType) NestedType() RestType {
	return nil
}

func (t *restFindByUUIDObjectType) TFAttrWithDiag() bool {
	return false
}

func (t *restFindByUUIDObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restFindByUUIDObjectType) TKHToTF(value string) string {
	return "types.StringPointerValue(" + value + ".GetUuid())"
}

func (t *restFindByUUIDObjectType) SDKTypeName() string {
	return t.nestedType.SDKTypeName()
}
