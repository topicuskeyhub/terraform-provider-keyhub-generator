// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

type restClassType struct {
	reachable         bool
	inReadOnlyContext bool
	suffix            string
	superClass        RestType
	realSuperClass    RestType
	name              string
	discriminator     string
	properties        []*RestProperty
	dsType            *restClassType
}

func NewRestClassType(realSuperClass RestType, superClass RestType, name string, discriminator string, inReadOnlyContext bool) *restClassType {
	ret := &restClassType{
		suffix:            "RS",
		realSuperClass:    realSuperClass,
		superClass:        superClass,
		name:              name,
		discriminator:     discriminator,
		inReadOnlyContext: inReadOnlyContext,
	}
	if inReadOnlyContext {
		ret.suffix = ret.suffix + "RO"
	}
	return ret
}

func (t *restClassType) Reachable() bool {
	return t.reachable
}

func (t *restClassType) MarkReachable() {
	if t.reachable {
		return
	}
	t.reachable = true
	if t.superClass != nil {
		t.superClass.MarkReachable()
	}
	for _, prop := range t.properties {
		prop.Type.MarkReachable()
	}
}

func (t *restClassType) Extends(typeName string) bool {
	return t.name == typeName || (t.realSuperClass != nil && t.realSuperClass.Extends(typeName))
}

func (t *restClassType) IsObject() bool {
	return true
}

func (t *restClassType) InReadOnlyContext() bool {
	return t.inReadOnlyContext
}

func (t *restClassType) ObjectAttrTypesName() string {
	return FirstCharToLower(t.name) + "AttrTypes"
}

func (t *restClassType) DataStructName() string {
	return FirstCharToLower(t.name) + "Data"
}

func (t *restClassType) APITypeName() string {
	return t.name
}

func (t *restClassType) APIDiscriminator() string {
	return t.discriminator
}

func (t *restClassType) GoTypeName() string {
	if t.InReadOnlyContext() {
		return FirstCharToUpper(t.name) + "RO"
	} else {
		return FirstCharToUpper(t.name)
	}
}

func (t *restClassType) SDKInterfaceTypeName() string {
	return t.SDKTypeName() + "able"
}

func (t *restClassType) SDKTypeName() string {
	return "keyhubmodel." + FirstCharToUpper(t.name)
}

func (t *restClassType) SDKTypeConstructor() string {
	return "keyhubmodel.New" + FirstCharToUpper(t.name) + "()"
}

func (t *restClassType) AllProperties() []*RestProperty {
	if t.superClass == nil {
		ret := make([]*RestProperty, len(t.properties))
		copy(ret, t.properties)
		return ret
	}
	super := t.superClass.AllProperties()
	sub := make([]*RestProperty, 0)
	for _, pt := range t.properties {
		found := false
		for _, ps := range super {
			if pt.Name == ps.Name {
				found = true
				break
			}
		}
		if !found {
			sub = append(sub, pt)
		}
	}
	return append(super, sub...)
}

func (t *restClassType) HasDirectUUIDProperty() bool {
	for _, prop := range t.properties {
		if prop.Name == "uuid" {
			return true
		}
	}
	if t.realSuperClass != nil {
		return t.realSuperClass.HasDirectUUIDProperty()
	}
	return false
}

func (t *restClassType) Suffix() string {
	return t.suffix
}

func (t *restClassType) DS() RestType {
	if t.dsType != nil {
		// break recursion
		return t.dsType
	}

	t.dsType = &restClassType{
		suffix:            "DS",
		name:              t.name,
		inReadOnlyContext: t.inReadOnlyContext,
	}
	if t.inReadOnlyContext {
		t.dsType.suffix = t.dsType.suffix + "RO"
	}
	if t.realSuperClass != nil {
		t.dsType.realSuperClass = t.realSuperClass.DS()
	}
	if t.superClass != nil {
		t.dsType.superClass = t.superClass.DS()
	}
	rsProperties := make([]*RestProperty, 0)
	for _, prop := range t.properties {
		if !prop.WriteOnly {
			rsProperties = append(rsProperties, prop.DS())
		}
	}
	t.dsType.properties = rsProperties
	return t.dsType
}
