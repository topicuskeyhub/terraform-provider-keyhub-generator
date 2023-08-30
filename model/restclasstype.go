package model

type restClassType struct {
	SuperClass RestType
	Name       string
	Properties []*RestProperty
}

func (t *restClassType) Extends(typeName string) bool {
	return t.Name == typeName || (t.SuperClass != nil && t.SuperClass.Extends(typeName))
}

func (t *restClassType) IsObject() bool {
	return true
}

func (t *restClassType) ObjectAttrTypesName() string {
	return firstCharToLower(t.Name) + "AttrTypes"
}

func (t *restClassType) DataStructName() string {
	return firstCharToLower(t.Name) + "Data"
}

func (t *restClassType) GoTypeName() string {
	return firstCharToUpper(t.Name)
}

func (t *restClassType) SDKTypeName() string {
	return "keyhubmodel." + firstCharToUpper(t.Name) + "able"
}

func (t *restClassType) SDKTypeConstructor() string {
	return "keyhubmodel.New" + firstCharToUpper(t.Name) + "()"
}

func (t *restClassType) AllProperties() []*RestProperty {
	if t.SuperClass == nil {
		ret := make([]*RestProperty, len(t.Properties))
		copy(ret, t.Properties)
		return ret
	}
	super := t.SuperClass.AllProperties()
	sub := make([]*RestProperty, 0)
	for _, pt := range t.Properties {
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
