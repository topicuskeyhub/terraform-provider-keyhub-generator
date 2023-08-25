package model

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type RestType interface {
	IsObject() bool
	ObjectAttrTypesName() string
	DataStructName() string
	GoTypeName() string
	SDKTypeName() string
	AllProperties() []*RestProperty
}

type RestProperty struct {
	Name     string
	Type     RestPropertyType
	Required bool
}

type RestPropertyType interface {
	PropertyNameSuffix() string
	TFName() string
	TFAttrType() string
	TFAttrWithDiag() bool
	TFAttrNeeded() bool
	NestedType() RestType
	TKHToTF(value string, list bool) string
	SDKTypeName(list bool) string
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

func firstCharToUpper(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToUpper(r)) + input[i:]
}
