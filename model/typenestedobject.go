package model

type restNestedObjectType struct {
	nestedType RestType
}

func NewNestedObjectType(nestedType RestType) RestPropertyType {
	return &restNestedObjectType{
		nestedType: nestedType,
	}
}

func (t *restNestedObjectType) PropertyNameSuffix() string {
	return ""
}

func (t *restNestedObjectType) TFName() string {
	return "types.Object"
}

func (t *restNestedObjectType) TFAttrType() string {
	return "types.ObjectType{AttrTypes: " + t.nestedType.ObjectAttrTypesName() + "}"
}

func (t *restNestedObjectType) NestedType() RestType {
	return t.nestedType
}

func (t *restNestedObjectType) TFAttrWithDiag() bool {
	return true
}

func (t *restNestedObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restNestedObjectType) TKHToTF(value string) string {
	return "tkhToTFObject" + t.nestedType.GoTypeName() + "(" + value + ")"
}

func (t *restNestedObjectType) SDKTypeName() string {
	return t.nestedType.SDKTypeName()
}
