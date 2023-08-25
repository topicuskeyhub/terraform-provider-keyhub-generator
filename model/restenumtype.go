package model

import "log"

type restEnumType struct {
	Name   string
	Values []string
}

func (t *restEnumType) IsObject() bool {
	return false
}

func (t *restEnumType) ObjectAttrTypesName() string {
	log.Fatalf("Enum type %s has no attributes", t.Name)
	return ""
}

func (t *restEnumType) DataStructName() string {
	log.Fatalf("Enum type %s has no attributes", t.Name)
	return ""
}

func (t *restEnumType) GoTypeName() string {
	log.Fatalf("Enum type %s has no attributes", t.Name)
	return ""
}

func (t *restEnumType) SDKTypeName() string {
	return "keyhubmodel." + firstCharToUpper(t.Name)
}

func (t *restEnumType) AllProperties() []*RestProperty {
	log.Fatalf("Enum type %s has no attributes", t.Name)
	return nil
}
