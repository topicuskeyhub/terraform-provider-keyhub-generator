package model

type restArrayType struct {
	itemType RestPropertyType
}

func NewRestArrayType(itemType RestPropertyType) RestPropertyType {
	return &restArrayType{
		itemType: itemType,
	}
}

func (t *restArrayType) PropertyNameSuffix() string {
	return ""
}

func (t *restArrayType) TFName() string {
	return "types.List"
}

func (t *restArrayType) TFAttrType() string {
	return "types.ListType{ElemType: " + t.itemType.TFAttrType() + "}"
}

func (t *restArrayType) TFValueType() string {
	return "basetypes.ListValue"
}

func (t *restArrayType) Complex() bool {
	return false
}

func (t *restArrayType) NestedType() RestType {
	return t.itemType.NestedType()
}

func (t *restArrayType) TFAttrWithDiag() bool {
	return true
}

func (t *restArrayType) TFAttrNeeded() bool {
	return true
}

func (t *restArrayType) TKHToTF(value string, list bool) string {
	sdkType := t.itemType.SDKTypeName(true)
	var body string
	if t.itemType.TFAttrWithDiag() {
		body = "            val, d := " + t.itemType.TKHToTF("tkh", true) + "\n" +
			"            diags.Append(d...)\n" +
			"            return val\n"
	} else {
		body = "            return " + t.itemType.TKHToTF("tkh", true) + "\n"
	}
	return "sliceToTF(elemType, " + value + ", func(tkh " + sdkType + ", diags *diag.Diagnostics) attr.Value {\n" +
		body +
		"        })"
}

func (t *restArrayType) SDKTypeName(list bool) string {
	return "[]" + t.itemType.SDKTypeName(true)
}

func (t *restArrayType) DSSchemaTemplate() string {
	return "data_source_schema_attr_array.go.tmpl"
}

func (t *restArrayType) DSSchemaTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"ElementType": t.itemType.TFAttrType(),
	}
}
