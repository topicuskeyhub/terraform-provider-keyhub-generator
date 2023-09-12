package model

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type RestProperty struct {
	Parent     RestType
	Name       string
	Type       RestPropertyType
	Required   bool
	dsProperty *RestProperty
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
	DSSchemaTemplateData() map[string]any
	RSSchemaTemplate() string
	RSSchemaTemplateData() map[string]any
	DS() RestPropertyType
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

func (p *RestProperty) DS() *RestProperty {
	if p.dsProperty != nil {
		// break recursion
		return p.dsProperty
	}

	p.dsProperty = &RestProperty{
		Parent:   p.Parent,
		Name:     p.Name,
		Required: p.Required,
	}
	p.dsProperty.Type = p.Type.DS()
	return p.dsProperty
}

func firstCharToLower(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToLower(r)) + input[i:]
}

func firstCharToUpper(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToUpper(r)) + input[i:]
}