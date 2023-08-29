package model

type restEnumPropertyType struct {
	enumType RestType
}

func NewEnumPropertyType(enumType RestType) RestPropertyType {
	return &restEnumPropertyType{
		enumType: enumType,
	}
}

func (t *restEnumPropertyType) PropertyNameSuffix() string {
	return ""
}

func (t *restEnumPropertyType) TFName() string {
	return "types.String"
}

func (t *restEnumPropertyType) TFAttrType() string {
	return "types.StringType"
}

func (t *restEnumPropertyType) TFValueType() string {
	return "basetypes.StringValue"
}

func (t *restEnumPropertyType) Complex() bool {
	return false
}

func (t *restEnumPropertyType) NestedType() RestType {
	return t.enumType
}

func (t *restEnumPropertyType) TFAttrWithDiag() bool {
	return false
}

func (t *restEnumPropertyType) TFAttrNeeded() bool {
	return false
}

func (t *restEnumPropertyType) TKHToTF(value string, list bool) string {
	if list {
		return "types.StringValue(" + value + ".String())"
	}
	return "stringerToTF(" + value + ")"
}

func (t *restEnumPropertyType) SDKTypeName(list bool) string {
	return t.enumType.SDKTypeName()
}

func (t *restEnumPropertyType) DSSchemaTemplate() string {
	return "data_source_schema_attr_simple.go.tmpl"
}

func (t *restEnumPropertyType) DSSchemaTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Type": "dsschema.StringAttribute",
	}
}
