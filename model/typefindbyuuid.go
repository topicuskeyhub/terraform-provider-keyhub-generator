package model

type restFindByUUIDObjectType struct {
}

func NewFindByUUIDObjectType() RestPropertyType {
	return &restFindByUUIDObjectType{}
}

func (t *restFindByUUIDObjectType) PropertyNameSuffix() string {
	return "Uuid"
}

func (t *restFindByUUIDObjectType) TFName() string {
	return "types.String"
}

func (t *restFindByUUIDObjectType) NestedType() *RestType {
	return nil
}
