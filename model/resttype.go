package model

type RestType interface {
	Extends(typeName string) bool
	IsObject() bool
	ObjectAttrTypesName() string
	DataStructName() string
	GoTypeName() string
	SDKTypeName() string
	SDKTypeConstructor() string
	AllProperties() []*RestProperty
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
