package model

import (
	"log"
)

type restFindByUUIDClassType struct {
	superClass   RestType
	name         string
	uuidProperty *RestProperty
	nestedType   *restClassType
}

func (t *restFindByUUIDClassType) Extends(typeName string) bool {
	return t.nestedType.Extends(typeName)
}

func (t *restFindByUUIDClassType) IsObject() bool {
	return t.nestedType.IsObject()
}

func (t *restFindByUUIDClassType) ObjectAttrTypesName() string {
	return t.nestedType.ObjectAttrTypesName()
}

func (t *restFindByUUIDClassType) DataStructName() string {
	return t.nestedType.DataStructName()
}

func (t *restFindByUUIDClassType) APITypeName() string {
	return t.nestedType.APITypeName()
}

func (t *restFindByUUIDClassType) GoTypeName() string {
	return t.nestedType.GoTypeName()
}

func (t *restFindByUUIDClassType) SDKTypeName() string {
	return t.nestedType.SDKTypeName()
}

func (t *restFindByUUIDClassType) SDKTypeConstructor() string {
	return t.nestedType.SDKTypeConstructor()
}

func (t *restFindByUUIDClassType) AllProperties() []*RestProperty {
	log.Printf("Type %s -> %s", t.name, t.uuidProperty.GoName())
	ret := make([]*RestProperty, 0)
	ret = append(ret, t.uuidProperty)
	ret = append(ret, t.nestedType.properties...)
	return ret
}

func (t *restFindByUUIDClassType) Suffix() string {
	return "RS"
}

func (t *restFindByUUIDClassType) DS() RestType {
	return t.nestedType.DS()
}
