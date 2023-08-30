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

func (t *restArrayType) ToTFAttrWithDiag() bool {
	return true
}

func (t *restArrayType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restArrayType) TFAttrNeeded() bool {
	return true
}

func (t *restArrayType) TKHToTF(value string, listItem bool) string {
	sdkType := t.itemType.SDKTypeName(true)
	var body string
	if t.itemType.ToTFAttrWithDiag() {
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

func (t *restArrayType) TFToTKH(value string, listItem bool) string {
	sdkType := t.itemType.SDKTypeName(true)
	var body string
	if t.itemType.ToTKHAttrWithDiag() {
		body = "            tkh, d := " + t.itemType.TFToTKH("val", true) + "\n" +
			"            diags.Append(d...)\n" +
			"            return tkh\n"
	} else {
		body = "            return " + t.itemType.TFToTKH("val", true) + "\n"
	}
	return "tfToSlice(" + value + ".(basetypes.ListValue), func(val attr.Value, diags *diag.Diagnostics) " + sdkType + " {\n" +
		body +
		"        })"
}

func (t *restArrayType) SDKTypeName(listItem bool) string {
	return "[]" + t.itemType.SDKTypeName(true)
}

func (t *restArrayType) SDKTypeConstructor() string {
	return "make([]" + t.itemType.SDKTypeName(true) + ", 0)"
}

func (t *restArrayType) DSSchemaTemplate() string {
	return "data_source_schema_attr_array.go.tmpl"
}

func (t *restArrayType) DSSchemaTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"ElementType": t.itemType.TFAttrType(),
	}
}
