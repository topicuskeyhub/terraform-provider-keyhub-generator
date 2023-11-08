// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

type restFindByUUIDClassType struct {
	reachable    bool
	superClass   RestType
	name         string
	uuidProperty *RestProperty
	nestedType   *restClassType
}

func NewRestFindByUUIDClassType(superClass RestType, name string, nestedType *restClassType) RestType {
	uuidType := &restFindByUUIDClassType{
		superClass: superClass,
		name:       name,
		nestedType: nestedType,
	}
	uuidType.uuidProperty = &RestProperty{
		Parent:   uuidType,
		Name:     "uuid",
		Required: true,
		Type:     NewFindBaseByUUIDObjectType(uuidType),
	}
	return uuidType
}

func (t *restFindByUUIDClassType) Reachable() bool {
	return t.reachable
}

func (t *restFindByUUIDClassType) MarkReachable() {
	if t.reachable {
		return
	}
	t.reachable = true
	t.superClass.MarkReachable()
	t.uuidProperty.Type.MarkReachable()
	t.nestedType.MarkReachable()
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

func (t *restFindByUUIDClassType) APIDiscriminator() string {
	return t.nestedType.APIDiscriminator()
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
	ret := make([]*RestProperty, 0)
	ret = append(ret, t.uuidProperty)
	ret = append(ret, t.nestedType.properties...)
	return ret
}

func (t *restFindByUUIDClassType) HasDirectUUIDProperty() bool {
	return false
}

func (t *restFindByUUIDClassType) Suffix() string {
	return "RS"
}

func (t *restFindByUUIDClassType) DS() RestType {
	return t.nestedType.DS()
}
