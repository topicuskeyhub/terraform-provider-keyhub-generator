package model

type restAdditionalType struct {
	names []string
}

func NewAdditionalType(names []string) RestPropertyType {
	return &restAdditionalType{
		names: names,
	}
}

func (t *restAdditionalType) PropertyNameSuffix() string {
	return ""
}

func (t *restAdditionalType) TFName() string {
	return "types.List"
}

func (t *restAdditionalType) TFAttrType() string {
	return "types.ListType{ElemType: types.StringType}"
}

func (t *restAdditionalType) TFValueType() string {
	return "basetypes.ListValue"
}

func (t *restAdditionalType) Complex() bool {
	return false
}

func (t *restAdditionalType) NestedType() RestType {
	return nil
}

func (t *restAdditionalType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restAdditionalType) ToTKHAttrWithDiag() bool {
	return false
}

func (t *restAdditionalType) ToTKHCustomCode() string {
	return ""
}

func (t *restAdditionalType) TFAttrNeeded() bool {
	return false
}

func (t *restAdditionalType) TKHToTF(value string, listItem bool) string {
	return ""
}

func (t *restAdditionalType) TFToTKH(value string, listItem bool) string {
	return ""
}

func (t *restAdditionalType) SDKTypeName(listItem bool) string {
	return "NONE"
}

func (t *restAdditionalType) SDKTypeConstructor() string {
	return "NONE"
}

func (t *restAdditionalType) DSSchemaTemplate() string {
	return "data_source_schema_attr_additional.go.tmpl"
}

func (t *restAdditionalType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Names": t.names,
	}
}

func (t *restAdditionalType) RSSchemaTemplate() string {
	return "NONE"
}

func (t *restAdditionalType) RSSchemaTemplateData() map[string]any {
	return map[string]any{}
}

func (t *restAdditionalType) DS() RestPropertyType {
	return nil
}
