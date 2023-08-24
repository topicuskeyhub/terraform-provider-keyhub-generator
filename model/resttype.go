package model

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/exp/maps"
)

type RestType struct {
	SuperClass *RestType
	Name       string
	Properties []*RestProperty
}

type RestProperty struct {
	Name     string
	Type     *RestPropertyType
	Required bool
}

type RestPropertyType interface {
	PropertyNameSuffix() string
	TFName() string
	NestedType() *RestType
}

func (t *RestType) DataStructName() string {
	return firstCharToUpper(t.Name) + "Data"
}

func (t *RestType) AllProperties() []*RestProperty {
	if t.SuperClass == nil {
		ret := make([]*RestProperty, len(t.Properties))
		copy(ret, t.Properties)
		return ret
	}
	return append(t.SuperClass.AllProperties(), t.Properties...)
}

func (t *RestType) AllRequiredTypes() []*RestType {
	types := make(map[string]*RestType)
	t.addAllRequiredTypes(types)
	return maps.Values(types)
}

func (t *RestType) addAllRequiredTypes(types map[string]*RestType) {
	types[t.Name] = t
	for _, prop := range t.Properties {
		if prop.Type != nil && (*prop.Type).NestedType() != nil {
			(*prop.Type).NestedType().addAllRequiredTypes(types)
		}
	}
}

func (p *RestProperty) internalName() string {
	if p.Type == nil {
		return p.Name
	}
	return p.Name + (*p.Type).PropertyNameSuffix()
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
	if p.Type == nil {
		return "UNKNOWN"
	}
	return (*p.Type).TFName()
}

func firstCharToUpper(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToUpper(r)) + input[i:]
}
