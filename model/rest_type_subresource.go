// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

type restSubresourceClassType struct {
	reachable         bool
	inReadOnlyContext bool
	name              string
	prefix            string
	nestedType        RestType
	dsType            *restSubresourceClassType
}

func NewRestSubresourceClassType(name string, prefix string, nestedType RestType, inReadOnlyContext bool) RestType {
	return &restSubresourceClassType{
		name:              name,
		prefix:            prefix,
		nestedType:        nestedType,
		inReadOnlyContext: inReadOnlyContext,
	}
}

func (t *restSubresourceClassType) Reachable() bool {
	return t.reachable
}

func (t *restSubresourceClassType) MarkReachable() {
	if t.reachable {
		return
	}
	t.reachable = true
	t.nestedType.MarkReachable()
}

func (t *restSubresourceClassType) Extends(typeName string) bool {
	return t.nestedType.Extends(typeName)
}

func (t *restSubresourceClassType) IsObject() bool {
	return t.nestedType.IsObject()
}

func (t *restSubresourceClassType) IsListOfFindByUuid() bool {
	return false
}

func (t *restSubresourceClassType) InReadOnlyContext() bool {
	return t.inReadOnlyContext
}

func (t *restSubresourceClassType) ObjectAttrTypesName() string {
	return FirstCharToLower(t.name) + "AttrTypes"
}

func (t *restSubresourceClassType) DataStructName() string {
	return FirstCharToLower(t.name) + "Data"
}

func (t *restSubresourceClassType) APITypeName() string {
	return t.nestedType.APITypeName()
}

func (t *restSubresourceClassType) APIDiscriminator() string {
	return t.nestedType.APIDiscriminator()
}

func (t *restSubresourceClassType) GoTypeName() string {
	return FirstCharToUpper(t.name)
}

func (t *restSubresourceClassType) SDKInterfaceTypeName() string {
	return t.nestedType.SDKInterfaceTypeName()
}

func (t *restSubresourceClassType) SDKTypeName() string {
	return t.nestedType.SDKTypeName()
}

func (t *restSubresourceClassType) SDKTypeConstructor() string {
	return t.nestedType.SDKTypeConstructor()
}

func (t *restSubresourceClassType) AllProperties() []*RestProperty {
	ret := make([]*RestProperty, 0)
	parent := &RestProperty{
		Parent:     t,
		Name:       t.prefix + "Uuid",
		Required:   true,
		WriteOnly:  false,
		Deprecated: false,
	}
	parent.Type = NewFindParentByUUIDObjectType()
	ret = append(ret, parent)
	ret = append(ret, t.nestedType.AllProperties()...)
	return ret
}

func (t *restSubresourceClassType) HasDirectUUIDProperty() bool {
	return t.nestedType.HasDirectUUIDProperty()
}

func (t *restSubresourceClassType) Suffix() string {
	return t.nestedType.Suffix()
}

func (t *restSubresourceClassType) DS() RestType {
	if t.dsType != nil {
		// break recursion
		return t.dsType
	}

	t.dsType = &restSubresourceClassType{
		name:              t.name,
		prefix:            t.prefix,
		nestedType:        t.nestedType.DS(),
		inReadOnlyContext: t.inReadOnlyContext,
	}
	return t.dsType
}
