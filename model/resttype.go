package model

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type RestType interface {
	Extends(typeName string) bool
	IsObject() bool
	ObjectAttrTypesName() string
	DataStructName() string
	GoTypeName() string
	SDKTypeName() string
	SDKTypeConstructor() string
	AllProperties() []*RestProperty
}

type RestProperty struct {
	Parent   RestType
	Name     string
	Type     RestPropertyType
	Required bool
}

type RestPropertyType interface {
	PropertyNameSuffix() string
	TFName() string
	TFValueType() string
	TFAttrType() string
	ToTFAttrWithDiag() bool
	ToTKHAttrWithDiag() bool
	TFAttrNeeded() bool
	Complex() bool
	NestedType() RestType
	TKHToTF(value string, listItem bool) string
	TFToTKH(value string, listItem bool) string
	SDKTypeName(listItem bool) string
	SDKTypeConstructor() string
	DSSchemaTemplate() string
	DSSchemaTemplateData() map[string]interface{}
}

func (p *RestProperty) internalName() string {
	return p.Name + p.Type.PropertyNameSuffix()
}

func (p *RestProperty) GoName() string {
	ret := firstCharToUpper(p.internalName())
	ret = strings.ReplaceAll(ret, "Uuid", "UUID")
	ret = strings.ReplaceAll(ret, "Id", "ID")
	ret = strings.ReplaceAll(ret, "Url", "URL")
	ret = strings.ReplaceAll(ret, "Tls", "TLS")
	return ret
}

func (p *RestProperty) TFName() string {
	ret := make([]rune, 0)
	for _, r := range p.internalName() {
		if unicode.IsUpper(r) {
			ret = append(ret, '_', unicode.ToLower(r))
		} else {
			ret = append(ret, r)
		}
	}
	return string(ret)
}

func (p *RestProperty) TFType() string {
	return p.Type.TFName()
}

func (p *RestProperty) TFAttrType() string {
	return p.Type.TFAttrType()
}

func (p *RestProperty) TKHToTF() string {
	return p.Type.TKHToTF("tkh.Get"+firstCharToUpper(p.Name)+"()", false)
}

func (p *RestProperty) TKHSetter() string {
	return "Set" + firstCharToUpper(p.Name)
}

func (p *RestProperty) TFToTKH() string {
	return p.Type.TFToTKH("objAttrs[\""+p.TFName()+"\"]", false)
}

func firstCharToLower(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToLower(r)) + input[i:]
}

func firstCharToUpper(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToUpper(r)) + input[i:]
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
