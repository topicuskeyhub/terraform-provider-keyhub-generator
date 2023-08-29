package model

type restNestedObjectType struct {
	property   *RestProperty
	nestedType RestType
}

func NewNestedObjectType(property *RestProperty, nestedType RestType) RestPropertyType {
	return &restNestedObjectType{
		property:   property,
		nestedType: nestedType,
	}
}

func (t *restNestedObjectType) PropertyNameSuffix() string {
	return ""
}

func (t *restNestedObjectType) TFName() string {
	return "types.Object"
}

func (t *restNestedObjectType) TFAttrType() string {
	return "types.ObjectType{AttrTypes: objectAttrsType" + t.nestedType.GoTypeName() + "(" + RecurseCutOff(t.property.Parent) + ")}"
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

func (t *restNestedObjectType) TFAttrWithDiag() bool {
	return true
}

func (t *restNestedObjectType) TFAttrNeeded() bool {
	return false
}

func (t *restNestedObjectType) TKHToTF(value string, list bool) string {
	return "tkhToTFObject" + t.nestedType.GoTypeName() + "(false, " + value + ")"
}

func (t *restNestedObjectType) SDKTypeName(list bool) string {
	return t.nestedType.SDKTypeName()
}

func (t *restNestedObjectType) DSSchemaTemplate() string {
	return "data_source_schema_attr_nestedobject.go.tmpl"
}

func (t *restNestedObjectType) DSSchemaTemplateData() map[string]interface{} {
	return make(map[string]interface{})
}
