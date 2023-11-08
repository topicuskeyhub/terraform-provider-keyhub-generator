// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

type RestType interface {
	Reachable() bool
	MarkReachable()
	Extends(typeName string) bool
	IsObject() bool
	ObjectAttrTypesName() string
	DataStructName() string
	APITypeName() string
	APIDiscriminator() string
	GoTypeName() string
	SDKTypeName() string
	SDKTypeConstructor() string
	AllProperties() []*RestProperty
	HasDirectUUIDProperty() bool
	Suffix() string
	DS() RestType
}

func RecurseCutOff(restType RestType) string {
	if AdditionalObjectsProperty(restType) != nil {
		return "false"
	}
	return "recurse"
}

func AdditionalObjectsProperty(restType RestType) *RestProperty {
	for _, curProperty := range restType.AllProperties() {
		if curProperty.Name == "additionalObjects" {
			return curProperty
		}
	}
	return nil
}

func AllDirectProperties(restType RestType) []*RestProperty {
	ret := make([]*RestProperty, 0)
	for _, curProperty := range restType.AllProperties() {
		if curProperty.Name != "additionalObjects" {
			ret = append(ret, curProperty)
		}
	}
	return ret
}
