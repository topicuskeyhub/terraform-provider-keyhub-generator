// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"fmt"

	"golang.org/x/exp/maps"
)

type restPolymorphicSubtype struct {
	property             *RestProperty
	superType            RestType
	nestedType           RestType
	rsSchemaTemplateBase map[string]any
}

func NewPolymorphicSubtype(property *RestProperty, superType RestType, nestedType RestType) RestPropertyType {
	return &restPolymorphicSubtype{
		property:   property,
		superType:  superType,
		nestedType: nestedType,
		rsSchemaTemplateBase: map[string]any{
			"Mode": "Optional",
		},
	}
}

func (t *restPolymorphicSubtype) MarkReachable() {
	t.superType.MarkReachable()
	t.nestedType.MarkReachable()
}

func (t *restPolymorphicSubtype) PropertyNameSuffix() string {
	return ""
}

func (t *restPolymorphicSubtype) FlattenMode() string {
	return "None"
}

func (t *restPolymorphicSubtype) OrderMode() string {
	return "Object"
}

func (t *restPolymorphicSubtype) TFName() string {
	return "types.Object"
}

func (t *restPolymorphicSubtype) TFAttrType(inAdditionalObjects bool) string {
	var recurseCutOff string
	if inAdditionalObjects {
		recurseCutOff = "false"
	} else {
		recurseCutOff = RecurseCutOff(t.property.Parent)
	}
	return "types.ObjectType{AttrTypes: objectAttrsType" +
		t.nestedType.Suffix() + t.nestedType.GoTypeName() + "(" + recurseCutOff + ")}"
}

func (t *restPolymorphicSubtype) TFValueType() string {
	return "basetypes.ObjectValue"
}

func (t *restPolymorphicSubtype) TFValidatorType() string {
	return ""
}

func (t *restPolymorphicSubtype) TFValidators() []string {
	return nil
}

func (t *restPolymorphicSubtype) Complex() bool {
	return true
}

func (t *restPolymorphicSubtype) RequiresReplace() bool {
	return false
}

func (t *restPolymorphicSubtype) NestedType() RestType {
	return t.nestedType
}

func (t *restPolymorphicSubtype) ToTFAttrWithDiag() bool {
	return true
}

func (t *restPolymorphicSubtype) ToTKHAttrWithDiag() bool {
	return true
}

func (t *restPolymorphicSubtype) TFAttrNeeded() bool {
	return false
}

func (t *restPolymorphicSubtype) TKHToTF(value string, listItem bool) string {
	return "tkhToTFObject" + t.nestedType.Suffix() + t.nestedType.GoTypeName() +
		"(" + RecurseCutOff(t.property.Parent) + ", " + value + ")"
}

func (t *restPolymorphicSubtype) TFToTKH(value string, listItem bool) string {
	return "tfObjectToTKH" + t.nestedType.Suffix() + t.nestedType.GoTypeName() +
		"(ctx, " + RecurseCutOff(t.property.Parent) + ", " + value + ".(basetypes.ObjectValue))"
}

func (t *restPolymorphicSubtype) TKHToTFGuard() string {
	return "tkhCast, _ := tkh.(" + t.nestedType.SDKTypeName() + ")\n"
}

func (t *restPolymorphicSubtype) TFToTKHGuard() string {
	return "if !objAttrs[\"" + t.property.TFName() + "\"].IsNull() "
}

func (t *restPolymorphicSubtype) TKHGetter(propertyName string) string {
	return "tkhCast"
}

func (t *restPolymorphicSubtype) ToTKHCustomCode(buildType RestType) string {
	typename := t.nestedType.GoTypeName()
	superTypename := t.superType.GoTypeName()
	return fmt.Sprintf("dtype := val.GetTypeEscaped()\n"+
		"(*val.(*keyhubmodel.%s)).%s = *tkh.(*keyhubmodel.%s)\n"+
		"val.SetTypeEscaped(dtype)\n"+
		"tkh = val", typename, superTypename, superTypename)
}

func (t *restPolymorphicSubtype) SDKTypeName(listItem bool) string {
	return t.nestedType.SDKTypeName()
}

func (t *restPolymorphicSubtype) SDKTypeConstructor() string {
	return t.nestedType.SDKTypeConstructor()
}

func (t *restPolymorphicSubtype) DSSchemaTemplate() string {
	return "data_source_schema_attr_nestedobject.go.tmpl"
}

func (t *restPolymorphicSubtype) DSSchemaTemplateData() map[string]any {
	return make(map[string]any)
}

func (t *restPolymorphicSubtype) RSSchemaTemplate() string {
	return "resource_schema_attr_nestedobject.go.tmpl"
}

func (t *restPolymorphicSubtype) RSSchemaTemplateData() map[string]any {
	ret := make(map[string]any)
	maps.Copy(ret, t.rsSchemaTemplateBase)
	return ret
}

func (t *restPolymorphicSubtype) DS() RestPropertyType {
	return NewNestedObjectType(t.property.DS(), t.nestedType.DS(), t.rsSchemaTemplateBase)
}
