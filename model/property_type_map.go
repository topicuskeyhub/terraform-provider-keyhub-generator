// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"strings"

	"golang.org/x/exp/maps"
)

type restMapType struct {
	name                 string
	itemType             RestPropertyType
	rsSchemaTemplateBase map[string]any
}

func NewRestMapType(name string, itemType RestPropertyType, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restMapType{
		name:                 name,
		itemType:             itemType,
		rsSchemaTemplateBase: rsSchemaTemplateBase,
	}
}

func (t *restMapType) ResolveRenderPropertyType() {
}

func (t *restMapType) MarkReachable() {
	t.itemType.MarkReachable()
}

func (t *restMapType) PropertyNameSuffix() string {
	return ""
}

func (t *restMapType) FlattenMode() string {
	return "None"
}

func (t *restMapType) OrderMode() string {
	return "Map"
}

func (t *restMapType) TFName() string {
	return "types.Map"
}

func (t *restMapType) TFAttrType(inAdditionalObjects bool) string {
	return "types.MapType{ElemType: " + t.itemType.TFAttrType(inAdditionalObjects) + "}"
}

func (t *restMapType) TFValueType() string {
	return "basetypes.MapValue"
}

func (t *restMapType) TFValueTypeCast() string {
	return "toMapValue"
}

func (t *restMapType) TFValidatorType() string {
	return "validator.Map"
}

func (t *restMapType) TFValidators() []string {
	if len(t.itemType.TFValidators()) == 0 {
		return nil
	}
	typename := t.itemType.TFValidatorType()
	typename = typename[strings.LastIndex(typename, ".")+1:]

	var sb strings.Builder
	sb.WriteString("listvalidator.Value")
	sb.WriteString(typename)
	sb.WriteString("sAre(\n")
	for _, validator := range t.itemType.TFValidators() {
		sb.WriteString(validator)
		sb.WriteString("\n")
	}
	sb.WriteString("),")
	return []string{sb.String()}
}

func (t *restMapType) Complex() bool {
	return t.itemType.Complex()
}

func (t *restMapType) RequiresReplace() bool {
	return false
}

func (t *restMapType) NestedType() RestType {
	return t.itemType.NestedType()
}

func (t *restMapType) ToTFAttrWithDiag() bool {
	return true
}

func (t *restMapType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restMapType) ToTKHCustomCode(buildType RestType) string {
	return ""
}

func (t *restMapType) TFAttrNeeded() bool {
	return true
}

func (t *restMapType) TKHToTF(value string, listItem bool) string {
	sdkType := t.itemType.SDKInterfaceTypeName(false)
	var body string
	if t.itemType.ToTFAttrWithDiag() {
		body = "            val, d := " + t.itemType.TKHToTF("tkh.("+sdkType+")", false) + "\n" +
			"            diags.Append(d...)\n" +
			"            return val\n"
	} else {
		body = "            return " + t.itemType.TKHToTF("tkh.("+sdkType+")", false) + "\n"
	}

	return "mapToTF(elemType, " + value + ".GetAdditionalData(), func(tkh any, diags *diag.Diagnostics) attr.Value {\n" +
		body +
		"        })"
}

func (t *restMapType) TFToTKH(planValue string, configValue string, listItem bool) string {
	var body string
	if t.itemType.ToTKHAttrWithDiag() {
		body = "            tkh, d := " + t.itemType.TFToTKH("planValue", "configValue", true) + "\n" +
			"            diags.Append(d...)\n" +
			"            return tkh\n"
	} else {
		body = "            return " + t.itemType.TFToTKH("planValue", "configValue", true) + "\n"
	}

	elementFunction := "func(planValue attr.Value, configValue attr.Value, diags *diag.Diagnostics) any {\n" +
		body +
		"        }"

	return "tfToMap(" + t.TFValueTypeCast() + "(" + planValue + ")," + t.TFValueTypeCast() + "(" + configValue + "), " + elementFunction + ", " + t.SDKTypeConstructor() + ")"
}

func (t *restMapType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restMapType) TKHToTFGuard() string {
	return ""
}

func (t *restMapType) TFToTKHGuard() string {
	return ""
}

func (t *restMapType) SDKInterfaceTypeName(listItem bool) string {
	return "keyhubmodel." + FirstCharToUpper(t.name) + "able"
}

func (t *restMapType) SDKTypeConstructor() string {
	return "keyhubmodel.New" + FirstCharToUpper(t.name) + "()"
}

func (t *restMapType) DSSchemaTemplate() string {
	return "data_source_schema_attr_map.go.tmpl"
}

func (t *restMapType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"ElementType": t.itemType.TFAttrType(false),
	}
}

func (t *restMapType) RSSchemaTemplate() string {
	return "resource_schema_attr_map.go.tmpl"
}

func (t *restMapType) RSSchemaTemplateData() map[string]any {
	ret := map[string]any{
		"ElementType": t.itemType.TFAttrType(false),
	}
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restMapType) DS() RestPropertyType {
	return NewRestMapType(t.name, t.itemType.DS(), t.rsSchemaTemplateBase)
}
