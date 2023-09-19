package model

import (
	"unicode"
)

type restPolymorphicBaseClassType struct {
	nestedType RestType
	subtypes   []RestType
	dsType     *restPolymorphicBaseClassType
}

func (t *restPolymorphicBaseClassType) Extends(typeName string) bool {
	return t.nestedType.Extends(typeName)
}

func (t *restPolymorphicBaseClassType) IsObject() bool {
	return t.nestedType.IsObject()
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
			Name:       FirstCharToLower(stripLowercasePrefix(subtype.APITypeName())),
			Required:   false,
			WriteOnly:  false,
			Deprecated: false,
		}
		prop.Type = NewPolymorphicSubtype(prop, t, subtype)
		ret = append(ret, prop)
	}
	return ret
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
		nestedType: t.nestedType.DS(),
	}
	t.dsType.subtypes = make([]RestType, 0)
	for _, subtype := range t.subtypes {
		t.dsType.subtypes = append(t.dsType.subtypes, subtype.DS())
	}
	return t.dsType
}

func stripLowercasePrefix(name string) string {
	firstUpper := 0
	for i, c := range name {
		if unicode.IsUpper(c) {
			firstUpper = i
			break
		}
	}
	return name[firstUpper:]
}
