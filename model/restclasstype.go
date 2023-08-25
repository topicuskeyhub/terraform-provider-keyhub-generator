package model

type restClassType struct {
	SuperClass RestType
	Name       string
	Properties []*RestProperty
}

func (t *restClassType) IsObject() bool {
	return true
}

func (t *restClassType) ObjectAttrTypesName() string {
	return t.Name + "AttrTypes"
}

func (t *restClassType) DataStructName() string {
	return t.Name + "Data"
}

func (t *restClassType) GoTypeName() string {
	return firstCharToUpper(t.Name)
}

func (t *restClassType) SDKTypeName() string {
	return "keyhubmodel." + firstCharToUpper(t.Name) + "able"
}

func (t *restClassType) AllProperties() []*RestProperty {
	if t.SuperClass == nil {
		ret := make([]*RestProperty, len(t.Properties))
		copy(ret, t.Properties)
		return ret
	}
	super := t.SuperClass.AllProperties()
	for _, p := range t.Properties {
		for is, ps := range super {
			if p.Name == ps.Name {
				super = append(super[:is], super[is+1:]...)
				break
			}
		}
	}
	return append(super, t.Properties...)
}
