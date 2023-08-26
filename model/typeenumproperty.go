package model

type restEnumPropertyType struct {
	enumType RestType
}

func NewEnumPropertyType(enumType RestType) RestPropertyType {
	return &restEnumPropertyType{
		enumType: enumType,
	}
}

func (t *restEnumPropertyType) PropertyNameSuffix() string {
	return ""
}

func (t *restEnumPropertyType) TFName() string {
	return "types.String"
}

func (t *restEnumPropertyType) TFAttrType() string {
	return "types.StringType"
}

func (t *restEnumPropertyType) NestedType() RestType {
	return t.enumType
}

func (t *restEnumPropertyType) TFAttrWithDiag() bool {
	return false
}

func (t *restEnumPropertyType) TFAttrNeeded() bool {
	return false
}

func (t *restEnumPropertyType) TKHToTF(value string, list bool) string {
	if list {
		return "types.StringValue(" + value + ".String())"
	}
	return "StringerToTF(" + value + ")"
}

func (t *restEnumPropertyType) SDKTypeName(list bool) string {
	return t.enumType.SDKTypeName()
}
