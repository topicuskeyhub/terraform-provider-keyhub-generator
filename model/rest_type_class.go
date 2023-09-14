package model

type restClassType struct {
	suffix        string
	superClass    RestType
	name          string
	discriminator string
	properties    []*RestProperty
	dsType        *restClassType
}

func (t *restClassType) Extends(typeName string) bool {
	return t.name == typeName || (t.superClass != nil && t.superClass.Extends(typeName))
}

func (t *restClassType) IsObject() bool {
	return true
}

func (t *restClassType) ObjectAttrTypesName() string {
	return firstCharToLower(t.name) + "AttrTypes"
}

func (t *restClassType) DataStructName() string {
	return firstCharToLower(t.name) + "Data"
}

func (t *restClassType) APITypeName() string {
	return t.name
}

func (t *restClassType) APIDiscriminator() string {
	return t.discriminator
}

func (t *restClassType) GoTypeName() string {
	return firstCharToUpper(t.name)
}

func (t *restClassType) SDKTypeName() string {
	return "keyhubmodel." + firstCharToUpper(t.name) + "able"
}

func (t *restClassType) SDKTypeConstructor() string {
	return "keyhubmodel.New" + firstCharToUpper(t.name) + "()"
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

func (t *restClassType) Suffix() string {
	return t.suffix
}

func (t *restClassType) DS() RestType {
	if t.dsType != nil {
		// break recursion
		return t.dsType
	}

	t.dsType = &restClassType{
		suffix: "DS",
		name:   t.name,
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
