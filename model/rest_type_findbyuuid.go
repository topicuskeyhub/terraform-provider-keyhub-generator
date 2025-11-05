// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

type restFindByUUIDClassType struct {
	reachable         bool
	inReadOnlyContext bool
	superClass        RestType
	name              string
	uuidProperty      *RestProperty
	nestedType        *restClassType
}

func NewRestFindByUUIDClassType(superClass RestType, name string, nestedType *restClassType, inReadOnlyContext bool) RestType {
	uuidType := &restFindByUUIDClassType{
		superClass:        superClass,
		name:              name,
		nestedType:        nestedType,
		inReadOnlyContext: inReadOnlyContext,
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

func (t *restFindByUUIDClassType) ResolveRenderPropertyType() RestType {
	return t
}

func (t *restFindByUUIDClassType) Extends(typeName string) bool {
	return t.nestedType.Extends(typeName)
}

func (t *restFindByUUIDClassType) IsObject() bool {
	return t.nestedType.IsObject()
}

func (t *restFindByUUIDClassType) IsListOfFindByUuid() bool {
	return false
}

func (t *restFindByUUIDClassType) InReadOnlyContext() bool {
	return t.inReadOnlyContext
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

func (t *restFindByUUIDClassType) SDKInterfaceTypeName() string {
	return t.nestedType.SDKInterfaceTypeName()
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
	if t.inReadOnlyContext {
		return "RSRO"
	} else {
		return "RS"
	}
}

func (t *restFindByUUIDClassType) DS() RestType {
	return t.nestedType.DS()
}
