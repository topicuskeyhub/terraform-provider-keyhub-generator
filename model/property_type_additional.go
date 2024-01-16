// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"strings"
)

type restAdditionalType struct {
	names []string
}

func NewAdditionalType(names []string) RestPropertyType {
	return &restAdditionalType{
		names: names,
	}
}

func (t *restAdditionalType) MarkReachable() {
}

func (t *restAdditionalType) PropertyNameSuffix() string {
	return ""
}

func (t *restAdditionalType) FlattenMode() string {
	return "None"
}

func (t *restAdditionalType) TFName() string {
	return "types.List"
}

func (t *restAdditionalType) TFAttrType(inAdditionalObjects bool) string {
	return "types.ListType{ElemType: types.StringType}"
}

func (t *restAdditionalType) TFValueType() string {
	return "basetypes.ListValue"
}

func (t *restAdditionalType) TFValidatorType() string {
	return "validator.List"
}

func (t *restAdditionalType) TFValidators() []string {
	var sb strings.Builder
	sb.WriteString("listvalidator.ValueStringsAre(stringvalidator.OneOf(\n")
	for _, name := range t.names {
		sb.WriteString(`"`)
		sb.WriteString(name)
		sb.WriteString(`",`)
	}
	sb.WriteString("\n)),")
	return []string{sb.String()}
}

func (t *restAdditionalType) Complex() bool {
	return false
}

func (t *restAdditionalType) RequiresReplace() bool {
	return false
}

func (t *restAdditionalType) NestedType() RestType {
	return nil
}

func (t *restAdditionalType) ToTFAttrWithDiag() bool {
	return false
}

func (t *restAdditionalType) ToTKHAttrWithDiag() bool {
	return false
}

func (t *restAdditionalType) ToTKHCustomCode() string {
	return ""
}

func (t *restAdditionalType) TFAttrNeeded() bool {
	return false
}

func (t *restAdditionalType) TKHToTF(value string, listItem bool) string {
	return "types.ListNull(types.StringType)"
}

func (t *restAdditionalType) TFToTKH(value string, listItem bool) string {
	return ""
}

func (t *restAdditionalType) TKHToTFGuard() string {
	return ""
}

func (t *restAdditionalType) TFToTKHGuard() string {
	return ""
}

func (t *restAdditionalType) TKHGetter(propertyName string) string {
	return ""
}

func (t *restAdditionalType) SDKTypeName(listItem bool) string {
	return "NONE"
}

func (t *restAdditionalType) SDKTypeConstructor() string {
	return "NONE"
}

func (t *restAdditionalType) DSSchemaTemplate() string {
	return "data_source_schema_attr_additional.go.tmpl"
}

func (t *restAdditionalType) DSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Names": t.names,
	}
}

func (t *restAdditionalType) RSSchemaTemplate() string {
	return "resource_schema_attr_additional.go.tmpl"
}

func (t *restAdditionalType) RSSchemaTemplateData() map[string]any {
	return map[string]any{
		"Names": t.names,
	}
}

func (t *restAdditionalType) DS() RestPropertyType {
	return NewAdditionalType(t.names)
}
