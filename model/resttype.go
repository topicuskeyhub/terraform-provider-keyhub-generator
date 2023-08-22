package model

type RestType struct {
	SuperClass *RestType
	Name       string
	Properties []*RestProperty
}

type RestProperty struct {
	Name     string
	Type     string
	Required bool
}
