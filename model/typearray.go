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

func (t *restArrayType) NestedType() RestType {
	return t.itemType.NestedType()
}

func (t *restArrayType) TFAttrWithDiag() bool {
	return true
}

func (t *restArrayType) TFAttrNeeded() bool {
	return true
}

func (t *restArrayType) TKHToTF(value string) string {
	sdkType := t.itemType.SDKTypeName()
	var body string
	if t.itemType.TFAttrWithDiag() {
		body = "            val, d := " + t.itemType.TKHToTF("tkh") + "\n" +
			"            diags.Append(d...)\n" +
			"            return val\n"
	} else {
		body = "            return " + t.itemType.TKHToTF("tkh") + "\n"
	}
	return "SliceToTF(elemType, " + value + ", func(tkh " + sdkType + ", diags *diag.Diagnostics) attr.Value {\n" +
		body +
		"        })"
}

func (t *restArrayType) SDKTypeName() string {
	return "[]" + t.itemType.SDKTypeName()
}
