package model

type restNestedObjectType struct {
	nestedType *RestType
}

func NewNestedObjectType(nestedType *RestType) RestPropertyType {
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

func (t *restNestedObjectType) NestedType() *RestType {
	return t.nestedType
}
