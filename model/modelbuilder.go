// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

var allSchemas openapi3.Schemas
var writableSubclassCounts map[string]int

func BuildModel(openapi *openapi3.T) map[string]map[bool]RestType {
	allSchemas = openapi.Components.Schemas
	collectWritableSubclassCounts()
	subresources := collectSubResources(openapi)
	ret := make(map[string]map[bool]RestType, 100)
	for name, schema := range allSchemas {
		if name == "RequestRange" {
			continue
		}
		getOrBuildTypeModel(ret, name, schema, nil, false)
	}
	for name, parents := range subresources {
		for _, parent := range parents {
			var nestedName string
			prefix := FirstCharToLower(StripLowercasePrefix(parent))
			if len(parents) == 1 {
				nestedName = "nested" + FirstCharToUpper(name)
			} else {
				nestedName = prefix + FirstCharToUpper(name)
			}
			getOrBuildTypeModel(ret, nestedName, allSchemas[name], &parentResourceInfo{
				typeName:     nestedName,
				prefix:       prefix,
				originalName: name,
			}, false)
		}
	}
	markReachable(ret)
	return ret
}

func markReachable(types map[string]map[bool]RestType) {
	types["authAccount"][false].DS().MarkReachable()
	types["certificateCertificate"][false].DS().MarkReachable()
	types["clientClientApplication"][false].DS().MarkReachable()
	types["directoryAccountDirectory"][false].DS().MarkReachable()
	types["groupGroup"][false].DS().MarkReachable()
	types["groupGroupClassification"][false].DS().MarkReachable()
	types["organizationOrganizationalUnit"][false].DS().MarkReachable()
	types["serviceaccountServiceAccount"][false].DS().MarkReachable()
	types["provisioningProvisionedSystem"][false].DS().MarkReachable()
	types["vaultVaultRecord"][false].DS().MarkReachable()
	types["webhookWebhook"][false].DS().MarkReachable()

	types["clientClientApplication"][false].MarkReachable()
	types["clientApplicationVaultVaultRecord"][false].MarkReachable()
	types["groupVaultVaultRecord"][false].MarkReachable()
	types["groupGroup"][false].MarkReachable()
	types["nestedProvisioningGroupOnSystem"][false].MarkReachable()
	types["serviceaccountServiceAccount"][false].MarkReachable()
}

func collectWritableSubclassCounts() {
	writableSubclassCounts = make(map[string]int)
	for _, schema := range allSchemas {
		ownSchema := findOwnTypeSchema(schema)
		writescopeObj, ok := ownSchema.Value.Extensions["x-tkh-writescope"]
		if !ok {
			continue
		}

		writescope := writescopeObj.(map[string]any)
		for name, val := range writescope {
			if val.(bool) {
				writableSubclassCounts[name] = writableSubclassCounts[name] + 1
			}
		}
	}
	for name, val := range writableSubclassCounts {
		if val == 1 || countSubclasses(name) < 2 || strings.HasSuffix(name, "Primer") {
			delete(writableSubclassCounts, name)
		}
	}
	for name := range writableSubclassCounts {
		if findPolymorphicBaseType(name) != nil {
			delete(writableSubclassCounts, name)
		}
	}
}

func findPolymorphicBaseType(name string) *string {
	checkName := name
	for {
		checkType, ok := allSchemas[checkName]
		if !ok {
			return nil
		}
		if checkType.Value.AllOf == nil || len(checkType.Value.AllOf) < 2 {
			return nil
		}
		checkName = refToName(checkType.Value.AllOf[0].Ref)
		if _, ok = writableSubclassCounts[checkName]; ok {
			return &checkName
		}
	}
}

func countSubclasses(typeName string) int {
	ret := 0
	for _, schema := range allSchemas {
		if schema.Value.AllOf != nil && len(schema.Value.AllOf) > 1 {
			superType := schema.Value.AllOf[0]
			if refToName(superType.Ref) == typeName {
				ret++
			}
		}
	}
	return ret
}

func collectSubResources(openapi *openapi3.T) map[string][]string {
	getpaths := make(map[string]string)
	stripId := regexp.MustCompile(`\{[^{]*\}`)
	for _, str := range openapi.Paths.InMatchingOrder() {
		path := openapi.Paths.Find(str)
		if strings.HasSuffix(str, "}") && path.Get != nil {
			for _, schema := range path.Get.Responses.Status(200).Value.Content {
				getpaths[stripId.ReplaceAllString(str, "#")] = refToName(schema.Schema.Ref)
				break
			}
		}
	}

	subresources := make(map[string][]string)
	for path, typeName := range getpaths {
		if strings.Count(path, "#") == 2 {
			parentResource := path[0:(strings.Index(path, "#") + 1)]
			subresources[typeName] = append(subresources[typeName], getpaths[parentResource])
		}
	}
	return subresources
}

type parentResourceInfo struct {
	typeName     string
	originalName string
	prefix       string
}

func getOrBuildTypeModel(types map[string]map[bool]RestType, name string, schema *openapi3.SchemaRef,
	parentResourceInfo *parentResourceInfo, inReadOnlyContext bool) RestType {
	if ret, ok := types[name][inReadOnlyContext]; ok {
		return ret
	}

	originalName := name
	if parentResourceInfo != nil {
		originalName = parentResourceInfo.originalName
	}
	var superType RestType
	var realSuperType RestType
	polymorphicBaseType := findPolymorphicBaseType(originalName)
	if schema != nil && schema.Value.AllOf != nil {
		superTypeName := refToName(schema.Value.AllOf[0].Ref)
		realSuperType = getOrBuildTypeModel(types, superTypeName, schema.Value.AllOf[0], nil, inReadOnlyContext)
	}
	if polymorphicBaseType == nil {
		superType = realSuperType
	}

	ownType := findOwnTypeSchema(schema)
	if ownType.Value.Type.Is("string") && len(ownType.Value.Enum) > 0 {
		return NewRestEnumType(originalName, ownType.Value.Enum)
	} else {
		discriminator := ""
		if discriminatorVal, ok := ownType.Value.Extensions["x-tkh-discriminator"]; ok {
			discriminator = discriminatorVal.(string)
		}
		classType := NewRestClassType(realSuperType, superType, originalName, discriminator, inReadOnlyContext)

		var ret RestType
		if isWritableWithUnwritableSuperClass(classType, ownType) {
			ret = NewRestFindByUUIDClassType(superType, originalName, classType, inReadOnlyContext)
		} else if _, ok := writableSubclassCounts[originalName]; ok {
			ret = NewRestPolymorphicBaseClassType(classType, inReadOnlyContext)
		} else {
			ret = classType
		}

		if parentResourceInfo != nil {
			ret = NewRestSubresourceClassType(name, parentResourceInfo.prefix, ret, inReadOnlyContext)
		}

		if existing, ok := types[name][inReadOnlyContext]; ok {
			return existing
		}
		if types[name] == nil {
			types[name] = make(map[bool]RestType)
		}
		types[name][inReadOnlyContext] = ret
		classType.properties = buildProperties(classType, originalName, ownType, types, inReadOnlyContext)

		if polymorphicBaseType != nil {
			found := false
			polyType := types[*polymorphicBaseType][inReadOnlyContext].(*restPolymorphicBaseClassType)

			for _, t := range polyType.subtypes {
				if t.APITypeName() == classType.APITypeName() {
					found = true
				}
			}
			if !found {
				polyType.subtypes = append(polyType.subtypes, classType)
			}

			// ensure other variant of every subtype exists
			getOrBuildTypeModel(types, name, schema, parentResourceInfo, !inReadOnlyContext)
		}
		return ret
	}
}

func findOwnTypeSchema(schema *openapi3.SchemaRef) *openapi3.SchemaRef {
	if schema.Value.AllOf != nil {
		return schema.Value.AllOf[1]
	}
	return schema
}

func isWritableWithUnwritableSuperClass(restType *restClassType, schema *openapi3.SchemaRef) bool {
	if restType.superClass == nil {
		return false
	}
	writescopeObj, ok := schema.Value.Extensions["x-tkh-writescope"]
	if !ok {
		return false
	}
	writescope := writescopeObj.(map[string]any)
	ownScope, ok := writescope[restType.name]
	if !ok {
		return false
	}
	superScope, ok := writescope[restType.superClass.APITypeName()]
	if !ok {
		return false
	}
	return ownScope.(bool) && !superScope.(bool)
}

func buildProperties(parent *restClassType, baseTypeName string, schema *openapi3.SchemaRef, types map[string]map[bool]RestType, inReadOnlyContext bool) []*RestProperty {
	required := schema.Value.Required
	ret := make([]*RestProperty, 0)
	var additionalObjectsProp *RestProperty = nil
	for name, property := range schema.Value.Properties {
		if skipProperty(baseTypeName, name) {
			continue
		}
		rsSchemaTemplateBase := buildRSSchemaTemplateBase(schema, name)
		if name == "type" {
			if baseTypeName == "RestLink" || baseTypeName == "authPermission" {
				name = "typeEscaped"
			} else {
				name = baseTypeName + "Type"
			}
		} else if name == "vendor" {
			name = "vendorEscaped"
		}

		restProperty := &RestProperty{
			Parent:     parent,
			Name:       name,
			Required:   slices.Contains(required, name),
			Deprecated: is(property, deprecated),
			WriteOnly:  is(property, writeOnly),
		}
		restProperty.Type = buildType(parent, baseTypeName, name, property, types, restProperty, rsSchemaTemplateBase, inReadOnlyContext)
		ret = append(ret, restProperty)
		if name == "additionalObjects" {
			additionalObjectsProp = restProperty
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})
	if additionalObjectsProp != nil {
		names := make([]string, 0)
		for _, prop := range additionalObjectsProp.Type.NestedType().AllProperties() {
			if !prop.WriteOnly {
				names = append(names, prop.Name)
			}
		}
		additionalProp := &RestProperty{
			Parent:     parent,
			Name:       "additional",
			Required:   false,
			WriteOnly:  false,
			Deprecated: false,
			Type:       NewAdditionalType(names),
		}
		ret = append([]*RestProperty{additionalProp}, ret...)
	}
	return ret
}

func skipProperty(baseTypeName string, propertyName string) bool {
	if propertyName == "$type" {
		return true
	}
	if propertyName == "additionalObjects" && baseTypeName == "authInternalAccount" {
		return true
	}
	if propertyName == "system" && baseTypeName == "provisioningGroupOnSystem" {
		return true
	}
	return false
}

func buildType(parentType *restClassType, baseTypeName string, propertyName string, ref *openapi3.SchemaRef, types map[string]map[bool]RestType, restProperty *RestProperty, rsSchemaTemplateBase map[string]any, inReadOnlyContext bool) RestPropertyType {
	schema := ref.Value

	log.Print("Building type for " + parentType.name + "." + propertyName + ", readonly: " + strconv.FormatBool(inReadOnlyContext))
	curInReadOnlyContext := (inReadOnlyContext || is(ref, readOnly))

	if len(schema.AllOf) > 0 {
		if ref.Ref == "" {
			ref = schema.AllOf[0]
		}
		schema = schema.AllOf[0].Value
	}
	if schema.Type.Is("array") {
		return NewRestArrayType(buildType(parentType, baseTypeName, propertyName, schema.Items, types, restProperty, rsSchemaTemplateBase, curInReadOnlyContext), schema.UniqueItems, rsSchemaTemplateBase)
	} else if schema.AdditionalProperties.Schema != nil {
		return NewRestMapType(baseTypeName+"_"+propertyName,
			buildType(parentType, baseTypeName, propertyName, schema.AdditionalProperties.Schema, types, restProperty, rsSchemaTemplateBase, curInReadOnlyContext),
			rsSchemaTemplateBase)
	}
	if ref.Ref != "" && schema.Type.Is("string") && len(schema.Enum) > 0 {
		enumName := refToName(ref.Ref)
		enum := getOrBuildTypeModel(types, enumName, ref, nil, curInReadOnlyContext)
		return NewEnumPropertyType(restProperty, enum, rsSchemaTemplateBase)
	}
	if schema.Type.Is("boolean") || schema.Type.Is("integer") || schema.Type.Is("string") {
		return NewRestSimpleType(restProperty, schema, rsSchemaTemplateBase)
	}
	if is(ref, object) {
		nestedTypeName := refToName(ref.Ref)
		if nestedTypeName == "" {
			nestedTypeName = baseTypeName + "_" + propertyName
		}

		nested := getOrBuildTypeModel(types, nestedTypeName, ref, nil, curInReadOnlyContext)
		ret := NewNestedObjectType(restProperty, nested, rsSchemaTemplateBase)
		if ref.Ref != "" && is(ref, withUUID) {
			if useFindByUUID(parentType, nested, curInReadOnlyContext) {
				ret = NewFindByUUIDObjectType(ret, rsSchemaTemplateBase)
			}
		}
		return ret
	}

	log.Fatalf("Cannot construct a type for (%v+)", ref)
	return nil
}

func useFindByUUID(parentType *restClassType, nested RestType, inReadOnlyContext bool) bool {
	ret := !inReadOnlyContext && (parentType.Extends("Linkable") || strings.HasSuffix(parentType.name, "_additionalObjects") || strings.HasSuffix(parentType.name, "LinkableWrapper") || strings.HasSuffix(parentType.name, "LinkableWrapperWithCount")) &&
		nested.Extends("Linkable") && nested.HasDirectUUIDProperty()

	return ret
}

func is(ref *openapi3.SchemaRef, check func(*openapi3.Schema) bool) bool {
	if check(ref.Value) {
		return true
	}
	for _, part := range ref.Value.AllOf {
		if is(part, check) {
			return true
		}
	}
	return false
}

func withUUID(schema *openapi3.Schema) bool {
	_, ok := schema.Properties["uuid"]
	return ok
}

func object(schema *openapi3.Schema) bool {
	return schema.Type.Is("object")
}

func deprecated(schema *openapi3.Schema) bool {
	return schema.Deprecated
}

func readOnly(schema *openapi3.Schema) bool {
	return schema.ReadOnly
}

func writeOnly(schema *openapi3.Schema) bool {
	return schema.WriteOnly
}

func refToName(ref string) string {
	return ref[strings.LastIndex(ref, "/")+1:]
}

/*
Computed:
Read only field, returned from backend

Computed_UseStateForUnknown
Immutable field, returned from backend

Computed_Optional (Not support for now)
Optional field with a default value determined by the backend if not provided

Optional
Optional field that can be left empty

Optional_Default
Optional or required field that has a default value

Required
Required field that does not have a default value
*/
func buildRSSchemaTemplateBase(ref *openapi3.SchemaRef, propertyName string) map[string]any {
	required := slices.Contains(ref.Value.Required, propertyName)
	property := ref.Value.Properties[propertyName]
	if property.Value.AllOf != nil && len(property.Value.AllOf) > 0 {
		property = property.Value.AllOf[1]
	}
	defaultVal := property.Value.Extensions["x-tkh-default"]
	readOnly := property.Value.ReadOnly
	immutable := property.Value.Extensions["x-tkh-immutable"] != nil && property.Value.Extensions["x-tkh-immutable"].(bool)
	createOnly := property.Value.Extensions["x-tkh-create-only"] != nil && property.Value.Extensions["x-tkh-create-only"].(bool)
	backendDefault := property.Value.Extensions["x-tkh-backend-determines-default"] != nil && property.Value.Extensions["x-tkh-backend-determines-default"].(bool)
	writeOnly := property.Value.WriteOnly

	if immutable && !createOnly {
		return map[string]any{
			"Mode": "Computed_UseStateForUnknown",
		}
	}
	if readOnly && !createOnly {
		return map[string]any{
			"Mode": "Computed",
		}
	}
	if defaultVal != nil {
		return map[string]any{
			"Mode":    "Optional_Default",
			"Default": defaultVal,
		}
	}
	if required {
		return map[string]any{
			"Mode": "Required",
		}
	}
	if backendDefault {
		return map[string]any{
			"Mode": "Optional_Computed",
		}
	}
	if writeOnly {
		return map[string]any{
			"Mode": "WriteOnly",
		}
	}
	return map[string]any{
		"Mode": "Optional",
	}
}

func FirstCharToLower(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToLower(r)) + input[i:]
}

func FirstCharToUpper(input string) string {
	r, i := utf8.DecodeRuneInString(input)
	return string(unicode.ToUpper(r)) + input[i:]
}

func StripLowercasePrefix(name string) string {
	firstUpper := 0
	for i, c := range name {
		if unicode.IsUpper(c) {
			firstUpper = i
			break
		}
	}
	return name[firstUpper:]
}
