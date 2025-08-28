// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"sort"
)

type restPolymorphicBaseClassType struct {
	reachable         bool
	inReadOnlyContext bool
	nestedType        RestType
	subtypes          []RestType
	dsType            *restPolymorphicBaseClassType
}

func NewRestPolymorphicBaseClassType(nestedType RestType, inReadOnlyContext bool) RestType {
	return &restPolymorphicBaseClassType{
		nestedType:        nestedType,
		subtypes:          make([]RestType, 0),
		inReadOnlyContext: inReadOnlyContext,
	}
}

func (t *restPolymorphicBaseClassType) Reachable() bool {
	return t.reachable
}

func (t *restPolymorphicBaseClassType) MarkReachable() {
	if t.reachable {
		return
	}
	t.reachable = true
	t.nestedType.MarkReachable()
	sort.Slice(t.subtypes, func(i, j int) bool {
		return t.subtypes[i].APITypeName() < t.subtypes[j].APITypeName()
	})

	for _, sub := range t.subtypes {
		sub.MarkReachable()
	}
}

func (t *restPolymorphicBaseClassType) Extends(typeName string) bool {
	return t.nestedType.Extends(typeName)
}

func (t *restPolymorphicBaseClassType) IsObject() bool {
	return t.nestedType.IsObject()
}

func (t *restPolymorphicBaseClassType) IsListOfFindByUuid() bool {
	return false
}

func (t *restPolymorphicBaseClassType) InReadOnlyContext() bool {
	return t.inReadOnlyContext
}

func (t *restPolymorphicBaseClassType) ObjectAttrTypesName() string {
	return t.nestedType.ObjectAttrTypesName()
}

func (t *restPolymorphicBaseClassType) DataStructName() string {
	return t.nestedType.DataStructName()
}

func (t *restPolymorphicBaseClassType) APITypeName() string {
	return t.nestedType.APITypeName()
}

func (t *restPolymorphicBaseClassType) APIDiscriminator() string {
	return t.nestedType.APIDiscriminator()
}

func (t *restPolymorphicBaseClassType) GoTypeName() string {
	return t.nestedType.GoTypeName()
}

func (t *restPolymorphicBaseClassType) SDKInterfaceTypeName() string {
	return t.nestedType.SDKInterfaceTypeName()
}

func (t *restPolymorphicBaseClassType) SDKTypeName() string {
	return t.nestedType.SDKTypeName()
}

func (t *restPolymorphicBaseClassType) SDKTypeConstructor() string {
	return t.nestedType.SDKTypeConstructor()
}

func (t *restPolymorphicBaseClassType) AllProperties() []*RestProperty {
	ret := make([]*RestProperty, 0)
	ret = append(ret, t.nestedType.AllProperties()...)
	for _, subtype := range t.subtypes {
		prop := &RestProperty{
			Parent:     t.nestedType,
			Name:       FirstCharToLower(StripLowercasePrefix(subtype.APITypeName())),
			Required:   false,
			WriteOnly:  false,
			Deprecated: false,
		}
		prop.Type = NewPolymorphicSubtype(prop, t, subtype)
		ret = append(ret, prop)
	}
	return ret
}

func (t *restPolymorphicBaseClassType) HasDirectUUIDProperty() bool {
	return t.nestedType.HasDirectUUIDProperty()
}

func (t *restPolymorphicBaseClassType) Suffix() string {
	return t.nestedType.Suffix()
}

func (t *restPolymorphicBaseClassType) DS() RestType {
	if t.dsType != nil {
		// break recursion
		return t.dsType
	}

	t.dsType = &restPolymorphicBaseClassType{
		nestedType:        t.nestedType.DS(),
		inReadOnlyContext: t.inReadOnlyContext,
	}

	t.dsType.subtypes = make([]RestType, 0)
	for _, subtype := range t.subtypes {
		t.dsType.subtypes = append(t.dsType.subtypes, subtype.DS())
	}
	return t.dsType
}
