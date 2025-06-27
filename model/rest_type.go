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
	SDKInterfaceTypeName() string
	SDKTypeName() string
	SDKTypeConstructor() string
	AllProperties() []*RestProperty
	HasDirectUUIDProperty() bool
	Suffix() string
	DS() RestType
}

type PropertyWithType struct {
	Property *RestProperty
	Type     RestType
}

func RecurseCutOff(restType RestType) string {
	if len(AdditionalObjectsProperties(restType)) > 0 {
		return "false"
	}
	return "recurse"
}

func AdditionalObjectsProperties(restType RestType) []*RestProperty {
	ret := make([]*RestProperty, 0)
	for _, curProperty := range restType.AllProperties() {
		if curProperty.Name == "additionalObjects" || curProperty.Name == "additional" {
			ret = append(ret, curProperty)
		}
	}
	return ret
}

func AllDirectProperties(restType RestType) []*RestProperty {
	ret := make([]*RestProperty, 0)
	for _, curProperty := range restType.AllProperties() {
		if curProperty.Name != "additionalObjects" && curProperty.Name != "additional" {
			ret = append(ret, curProperty)
		}
	}
	return ret
}

func IdentifyingProperties(restType RestType) []*RestProperty {
	ret := make([]*RestProperty, 0)
	notComputed := make([]*RestProperty, 0)
	for _, p := range AllDirectProperties(restType) {
		if !p.Type.Complex() {
			if p.IsRequired() {
				ret = append(ret, p)
			} else if p.IsNotComputed() {
				notComputed = append(notComputed, p)
			}
		}
	}
	ret = append(ret, notComputed...)
	return ret
}

func ItemsProperty(properties []*RestProperty) *RestProperty {
	for _, p := range properties {
		if p.Name == "items" {
			return p
		}
	}
	return nil
}

func ToPropertyWithType(property *RestProperty, restType RestType) PropertyWithType {
	return PropertyWithType{
		Property: property,
		Type:     restType,
	}
}
