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

func (t *restArrayType) OrderMode() string {
	if t.setCollection || !t.Complex() {
		return "None"
	}
	return "List"
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

func (t *restArrayType) TFValueTypeCast() string {
	if t.setCollection {
		return "toSetValue"
	} else {
		return "toListValue"
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

func (t *restArrayType) RequiresReplace() bool {
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

func (t *restArrayType) ToTKHCustomCode(baseType RestType) string {
	return ""
}

func (t *restArrayType) TFAttrNeeded() bool {
	return true
}

func (t *restArrayType) TKHToTF(value string, listItem bool) string {
	sdkType := t.itemType.SDKInterfaceTypeName(true)
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

func (t *restArrayType) TFToTKH(planValue string, configValue string, listItem bool) string {
	sdkType := t.itemType.SDKInterfaceTypeName(true)
	var body string
	if t.itemType.ToTKHAttrWithDiag() {
		body = "            tkh, d := " + t.itemType.TFToTKH("planValue", "configValue", true) + "\n" +
			"            diags.Append(d...)\n" +
			"            return tkh\n"
	} else {
		body = "            return " + t.itemType.TFToTKH("planValue", "configValue", true) + "\n"
	}
	elementFunction := "func(planValue attr.Value, configValue attr.Value, diags *diag.Diagnostics) " + sdkType + " {\n" +
		body +
		"        }"

	var collectionFunctionName string
	if t.setCollection {
		collectionFunctionName = "tfToSliceSet"
	} else {
		collectionFunctionName = "tfToSliceListBinary"
	}

	// basically a foreach construction
	return collectionFunctionName + "(" + t.TFValueTypeCast() + "(" + planValue + ")," + t.TFValueTypeCast() + "(" + configValue + "), " + elementFunction + ")"
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

func (t *restArrayType) SDKInterfaceTypeName(listItem bool) string {
	return "[]" + t.itemType.SDKInterfaceTypeName(true)
}

func (t *restArrayType) SDKTypeConstructor() string {
	return "make([]" + t.itemType.SDKInterfaceTypeName(true) + ", 0)"
}

func (t *restArrayType) DSSchemaTemplate() string {
	return "data_source_schema_attr_array.go.tmpl"
}

func (t *restArrayType) fillSchemaTemplateBase() map[string]any {
	ret := map[string]any{
		"ElementType": t.itemType.TFAttrType(false),
	}
	if t.setCollection {
		if t.itemType.Complex() {
			ret["SchemaAttributeType"] = "SetNestedAttribute"
		} else {
			ret["SchemaAttributeType"] = "SetAttribute"
		}
		ret["StateForUnknown"] = "[]planmodifier.Set{setplanmodifier.UseStateForUnknown()}"
	} else {
		if t.itemType.Complex() {
			ret["SchemaAttributeType"] = "ListNestedAttribute"
		} else {
			ret["SchemaAttributeType"] = "ListAttribute"
		}
		ret["StateForUnknown"] = "[]planmodifier.List{listplanmodifier.UseStateForUnknown()}"
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restArrayType) DSSchemaTemplateData() map[string]any {
	return t.fillSchemaTemplateBase()
}

func (t *restArrayType) RSSchemaTemplate() string {
	return "resource_schema_attr_array.go.tmpl"
}

func (t *restArrayType) RSSchemaTemplateData() map[string]any {
	return t.fillSchemaTemplateBase()
}

func (t *restArrayType) DS() RestPropertyType {
	return NewRestArrayType(t.itemType.DS(), t.setCollection, t.rsSchemaTemplateBase)
}
