package model

import "golang.org/x/exp/maps"

type restArrayType struct {
	itemType             RestPropertyType
	rsSchemaTemplateBase map[string]any
}

func NewRestArrayType(itemType RestPropertyType, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restArrayType{
		itemType:             itemType,
		rsSchemaTemplateBase: rsSchemaTemplateBase,
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

func (t *restArrayType) TFValidatorType() string {
	return ""
}

func (t *restArrayType) TFValidators() []string {
	return nil
}

func (t *restArrayType) Complex() bool {
	return t.itemType.Complex()
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

func (t *restArrayType) ToTKHCustomCode() string {
	return ""
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

func (t *restArrayType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restArrayType) TKHToTFGuard() string {
	return ""
}

func (t *restArrayType) TFToTKHGuard() string {
	return ""
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

func (t *restArrayType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"ElementType": t.itemType.TFAttrType(),
	}
}

func (t *restArrayType) RSSchemaTemplate() string {
	return "resource_schema_attr_array.go.tmpl"
}

func (t *restArrayType) RSSchemaTemplateData() map[string]any {
	ret := map[string]any{
		"ElementType": t.itemType.TFAttrType(),
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restArrayType) DS() RestPropertyType {
	return NewRestArrayType(t.itemType.DS(), t.rsSchemaTemplateBase)
}
