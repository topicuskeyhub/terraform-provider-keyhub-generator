package model

type restArrayType struct {
	itemType *RestPropertyType
}

func NewRestArrayType(itemType *RestPropertyType) RestPropertyType {
	return &restArrayType{
		itemType: itemType,
	}
}

func (t *restArrayType) PropertyNameSuffix() string {
	return ""
}

func (t *restArrayType) TFName() string {
	return "types.List"
}

func (t *restArrayType) NestedType() *RestType {
	if t.itemType == nil {
		return nil
	}
	return (*t.itemType).NestedType()
}
