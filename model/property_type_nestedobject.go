package model

import "golang.org/x/exp/maps"

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

func (t *restNestedObjectType) PropertyNameSuffix() string {
	return ""
}

func (t *restNestedObjectType) TFName() string {
	return "types.Object"
}

func (t *restNestedObjectType) TFAttrType() string {
	return "types.ObjectType{AttrTypes: objectAttrsType" +
		t.nestedType.Suffix() + t.nestedType.GoTypeName() + "(" + RecurseCutOff(t.property.Parent) + ")}"
}

func (t *restNestedObjectType) TFValueType() string {
	return "basetypes.ObjectValue"
}

func (t *restNestedObjectType) Complex() bool {
	return true
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

func (t *restNestedObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restNestedObjectType) TKHToTF(value string, listItem bool) string {
	return "tkhToTFObject" + t.nestedType.Suffix() + t.nestedType.GoTypeName() +
		"(" + RecurseCutOff(t.property.Parent) + ", " + value + ")"
}

func (t *restNestedObjectType) TFToTKH(value string, listItem bool) string {
	return "tfObjectToTKH" + t.nestedType.Suffix() + t.nestedType.GoTypeName() +
		"(ctx, " + RecurseCutOff(t.property.Parent) + ", " + value + ".(basetypes.ObjectValue))"
}

func (t *restNestedObjectType) SDKTypeName(listItem bool) string {
	return t.nestedType.SDKTypeName()
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