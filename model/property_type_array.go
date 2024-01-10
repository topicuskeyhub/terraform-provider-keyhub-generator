// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"strings"

	"golang.org/x/exp/maps"
)

type restArrayType struct {
	itemType             RestPropertyType
	setCollection        bool
	rsSchemaTemplateBase map[string]any
}

func NewRestArrayType(itemType RestPropertyType, setCollection bool, rsSchemaTemplateBase map[string]any) RestPropertyType {
	if setCollection {
		rsSchemaTemplateBase["SchemaAttributeType"] = "SetAttribute"
	} else {
		rsSchemaTemplateBase["SchemaAttributeType"] = "ListAttribute"
	}
	return &restArrayType{
		itemType:             itemType,
		setCollection:        setCollection,
		rsSchemaTemplateBase: rsSchemaTemplateBase,
	}
}

func (t *restArrayType) MarkReachable() {
	t.itemType.MarkReachable()
}

func (t *restArrayType) PropertyNameSuffix() string {
	return ""
}

func (t *restArrayType) FlattenMode() string {
	return "None"
}

func (t *restArrayType) TFName() string {
	if t.setCollection {
		return "types.Set"
	} else {
		return "types.List"
	}
}

func (t *restArrayType) TFAttrType(inAdditionalObjects bool) string {
	var structName string
	if t.setCollection {
		structName = "types.SetType"
	} else {
		structName = "types.ListType"
	}
	return structName + "{ElemType: " + t.itemType.TFAttrType(inAdditionalObjects) + "}"
}

func (t *restArrayType) TFValueType() string {
	if t.setCollection {
		return "basetypes.SetValue"
	} else {
		return "basetypes.ListValue"
	}
}

func (t *restArrayType) TFValidatorType() string {
	if t.setCollection {
		return "validator.Set"
	} else {
		return "validator.List"
	}
}

func (t *restArrayType) TFValidators() []string {
	if len(t.itemType.TFValidators()) == 0 {
		return nil
	}
	typename := t.itemType.TFValidatorType()
	typename = typename[strings.LastIndex(typename, ".")+1:]

	var sb strings.Builder
	if t.setCollection {
		sb.WriteString("setvalidator.Value")
	} else {
		sb.WriteString("listvalidator.Value")
	}
	sb.WriteString(typename)
	sb.WriteString("sAre(\n")
	for _, validator := range t.itemType.TFValidators() {
		sb.WriteString(validator)
		sb.WriteString("\n")
	}
	sb.WriteString("),")
	return []string{sb.String()}
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
	var functionName string
	if t.setCollection {
		functionName = "sliceToTFSet"
	} else {
		functionName = "sliceToTFList"
	}
	return functionName + "(elemType, " + value + ", func(tkh " + sdkType + ", diags *diag.Diagnostics) attr.Value {\n" +
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
	var functionName string
	if t.setCollection {
		functionName = "tfToSliceSet"
	} else {
		functionName = "tfToSliceList"
	}
	return functionName + "(" + value + ".(" + t.TFValueType() + "), func(val attr.Value, diags *diag.Diagnostics) " + sdkType + " {\n" +
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
	ret := map[string]any{
		"ElementType": t.itemType.TFAttrType(false),
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restArrayType) RSSchemaTemplate() string {
	return "resource_schema_attr_array.go.tmpl"
}

func (t *restArrayType) RSSchemaTemplateData() map[string]any {
	ret := map[string]any{
		"ElementType": t.itemType.TFAttrType(false),
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restArrayType) DS() RestPropertyType {
	return NewRestArrayType(t.itemType.DS(), t.setCollection, t.rsSchemaTemplateBase)
}
