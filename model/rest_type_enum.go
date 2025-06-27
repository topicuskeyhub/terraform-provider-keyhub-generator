// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import "log"

type restEnumType struct {
	reachable bool
	suffix    string
	name      string
	values    []any
}

func NewRestEnumType(name string, values []any) RestType {
	return &restEnumType{
		suffix: "RS",
		name:   name,
		values: values,
	}
}

func (t *restEnumType) Reachable() bool {
	return t.reachable
}

func (t *restEnumType) MarkReachable() {
	t.reachable = true
}

func (t *restEnumType) Extends(typeName string) bool {
	return false
}

func (t *restEnumType) IsObject() bool {
	return false
}

func (t *restEnumType) ObjectAttrTypesName() string {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return ""
}

func (t *restEnumType) DataStructName() string {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return ""
}

func (t *restEnumType) APITypeName() string {
	return t.name
}

func (t *restEnumType) APIDiscriminator() string {
	return ""
}

func (t *restEnumType) GoTypeName() string {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return ""
}

func (t *restEnumType) SDKInterfaceTypeName() string {
	return t.SDKTypeName()
}

func (t *restEnumType) SDKTypeName() string {
	return "keyhubmodel." + FirstCharToUpper(t.name)
}

func (t *restEnumType) SDKTypeConstructor() string {
	return "keyhubmodel.Parse" + FirstCharToUpper(t.name)
}

func (t *restEnumType) AllProperties() []*RestProperty {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return nil
}

func (t *restEnumType) HasDirectUUIDProperty() bool {
	return false
}

func (t *restEnumType) Suffix() string {
	return t.suffix
}

func (t *restEnumType) DS() RestType {
	return &restEnumType{
		suffix: "DS",
		name:   t.name,
		values: t.values,
	}
}
