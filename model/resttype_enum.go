package model

import "log"

type restEnumType struct {
	suffix string
	name   string
	values []string
}

func (t *restEnumType) Extends(typeName string) bool {
	return false
}

func (t *restEnumType) IsObject() bool {
	return false
}

func (t *restEnumType) ObjectAttrTypesName() string {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return ""
}

func (t *restEnumType) DataStructName() string {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return ""
}

func (t *restEnumType) GoTypeName() string {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return ""
}

func (t *restEnumType) SDKTypeName() string {
	return "keyhubmodel." + firstCharToUpper(t.name)
}

func (t *restEnumType) SDKTypeConstructor() string {
	return "keyhubmodel.Parse" + firstCharToUpper(t.name)
}

func (t *restEnumType) AllProperties() []*RestProperty {
	log.Fatalf("Enum type %s has no attributes", t.name)
	return nil
}

func (t *restEnumType) Suffix() string {
	return t.suffix
}

func (t *restEnumType) DS() RestType {
	return &restEnumType{
		suffix: "DS",
		name:   t.name,
		values: t.values,
	}
}
