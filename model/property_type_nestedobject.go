// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"strings"

	"golang.org/x/exp/maps"
)

type restNestedObjectType struct {
	property             *RestProperty
	nestedType           RestType
	rsSchemaTemplateBase map[string]any
}

func NewNestedObjectType(property *RestProperty, nestedType RestType, rsSchemaTemplateBase map[string]any) RestPropertyType {
	return &restNestedObjectType{
		property:             property,
		nestedType:           nestedType,
		rsSchemaTemplateBase: rsSchemaTemplateBase,
	}
}

func (t *restNestedObjectType) MarkReachable() {
	t.nestedType.MarkReachable()
}

func (t *restNestedObjectType) PropertyNameSuffix() string {
	return ""
}

func (t *restNestedObjectType) FlattenMode() string {
	if t.property.TFName() == "additional_objects" {
		return "AdditionalObjects"
	}
	if strings.HasSuffix(t.nestedType.APITypeName(), "LinkableWrapper") ||
		strings.HasSuffix(t.nestedType.APITypeName(), "LinkableWrapperWithCount") {
		nestedProps := t.nestedType.AllProperties()
		if len(nestedProps) == 1 && nestedProps[0].TFName() == "items" {
			return "ItemsList"
		}
		if len(nestedProps) == 2 &&
			(nestedProps[0].TFName() == "items" || nestedProps[1].TFName() == "items") {
			return "ItemsList"
		}
	}
	return "None"
}

func (t *restNestedObjectType) OrderMode() string {
	if t.property.TFName() == "additional_objects" {
		return "AdditionalObjects"
	}
	if strings.HasSuffix(t.nestedType.APITypeName(), "LinkableWrapper") ||
		strings.HasSuffix(t.nestedType.APITypeName(), "LinkableWrapperWithCount") {
		nestedProps := t.nestedType.AllProperties()
		if len(nestedProps) == 1 && nestedProps[0].TFName() == "items" {
			return "ItemsList"
		}
		if len(nestedProps) == 2 &&
			(nestedProps[0].TFName() == "items" || nestedProps[1].TFName() == "items") {
			return "ItemsList"
		}
	}
	return "Object"
}

func (t *restNestedObjectType) TFName() string {
	if t.FlattenMode() == "ItemsList" {
		return "types.List"
	}
	return "types.Object"
}

func (t *restNestedObjectType) TFAttrType(inAdditionalObjects bool) string {
	var recurseCutOff string
	if inAdditionalObjects {
		recurseCutOff = "false"
	} else {
		recurseCutOff = RecurseCutOff(t.property.Parent)
	}
	nestedAttrType := "objectAttrsType" + t.nestedType.Suffix() + t.nestedType.GoTypeName() + "(" + recurseCutOff + ")"
	if t.FlattenMode() == "ItemsList" {
		return nestedAttrType + `["items"]`
	}
	return "types.ObjectType{AttrTypes: " + nestedAttrType + "}"
}

func (t *restNestedObjectType) TFValueType() string {
	if t.FlattenMode() == "ItemsList" {
		return "basetypes.ListValue"
	}
	return "basetypes.ObjectValue"
}

func (t *restNestedObjectType) TFValidatorType() string {
	return ""
}

func (t *restNestedObjectType) TFValidators() []string {
	return nil
}

func (t *restNestedObjectType) Complex() bool {
	return true
}

func (t *restNestedObjectType) RequiresReplace() bool {
	return false
}

func (t *restNestedObjectType) NestedType() RestType {
	return t.nestedType
}

func (t *restNestedObjectType) ToTFAttrWithDiag() bool {
	return true
}

func (t *restNestedObjectType) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restNestedObjectType) ToTKHCustomCode(buildType RestType) string {
	return ""
}

func (t *restNestedObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restNestedObjectType) TKHToTF(value string, listItem bool) string {
	return "tkhToTFObject" + t.nestedType.Suffix() + t.nestedType.GoTypeName() +
		"(" + RecurseCutOff(t.property.Parent) + ", " + value + ")"
}

func (t *restNestedObjectType) TFToTKH(planValue string, configValue string, listItem bool) string {
	var tfPlanVal string
	var tfConfigVal string
	if t.FlattenMode() == "AdditionalObjects" {
		tfPlanVal = "planValues"
		tfConfigVal = "configValues"
	} else if t.FlattenMode() == "ItemsList" {
		tfPlanVal = `toItemsList(ctx, planAttrValues["` + t.property.TFName() + `"])`
		tfConfigVal = `toItemsList(ctx, configAttrValues["` + t.property.TFName() + `"])`
	} else {
		tfPlanVal = "toObjectValue(" + planValue + ")"
		tfConfigVal = "toObjectValue(" + configValue + ")"
	}

	return "tfObjectToTKH" + t.nestedType.Suffix() + t.nestedType.GoTypeName() +
		"(ctx, " + RecurseCutOff(t.property.Parent) + ", " + tfPlanVal + ", " + tfConfigVal + ")"
}

func (t *restNestedObjectType) TKHToTFGuard() string {
	return ""
}

func (t *restNestedObjectType) TFToTKHGuard() string {
	return ""
}

func (t *restNestedObjectType) TKHGetter(propertyName string) string {
	return "tkh.Get" + FirstCharToUpper(propertyName) + "()"
}

func (t *restNestedObjectType) SDKInterfaceTypeName(listItem bool) string {
	return t.nestedType.SDKInterfaceTypeName()
}

func (t *restNestedObjectType) SDKTypeConstructor() string {
	return t.nestedType.SDKTypeConstructor()
}

func (t *restNestedObjectType) DSSchemaTemplate() string {
	return "data_source_schema_attr_nestedobject.go.tmpl"
}

func (t *restNestedObjectType) DSSchemaTemplateData() map[string]any {
	return make(map[string]any)
}

func (t *restNestedObjectType) RSSchemaTemplate() string {
	return "resource_schema_attr_nestedobject.go.tmpl"
}

func (t *restNestedObjectType) RSSchemaTemplateData() map[string]any {
	ret := make(map[string]any)
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restNestedObjectType) DS() RestPropertyType {
	return NewNestedObjectType(t.property.DS(), t.nestedType.DS(), t.rsSchemaTemplateBase)
}
