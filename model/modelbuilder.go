package model

import (
	"log"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
)

var allSchemas openapi3.Schemas
var writableSubclassCounts map[string]int

func BuildModel(openapi *openapi3.T) map[string]RestType {
	allSchemas = openapi.Components.Schemas
	collectWritableSubclassCounts()
	subresources := collectSubResources(openapi)
	ret := make(map[string]RestType, 100)
	for name, schema := range allSchemas {
		if name == "RequestRange" {
			continue
		}
		getOrBuildTypeModel(ret, name, schema, nil)
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
			})
		}
	}
	return ret
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
		if val == 1 || countSubclasses(name) < 2 {
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
	for str, path := range openapi.Paths {
		if strings.HasSuffix(str, "}") {
			for _, schema := range path.Get.Responses["200"].Value.Content {
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

func getOrBuildTypeModel(types map[string]RestType, name string, schema *openapi3.SchemaRef,
	parentResourceInfo *parentResourceInfo) RestType {
	if ret, ok := types[name]; ok {
		return ret
	}

	originalName := name
	if parentResourceInfo != nil {
		originalName = parentResourceInfo.originalName
	}
	var superType RestType
	polymorphicBaseType := findPolymorphicBaseType(originalName)
	if schema != nil && schema.Value.AllOf != nil {
		superTypeName := refToName(schema.Value.AllOf[0].Ref)
		superType = getOrBuildTypeModel(types, superTypeName, schema.Value.AllOf[0], nil)
	}
	if polymorphicBaseType != nil {
		superType = nil
	}

	ownType := findOwnTypeSchema(schema)
	if ownType.Value.Type == "string" && len(ownType.Value.Enum) > 0 {
		ret := &restEnumType{
			suffix: "RS",
			name:   originalName,
		}
		return ret
	} else {
		discriminator := ""
		if discriminatorVal, ok := ownType.Value.Extensions["x-tkh-discriminator"]; ok {
			discriminator = discriminatorVal.(string)
		}
		classType := &restClassType{
			suffix:        "RS",
			superClass:    superType,
			discriminator: discriminator,
			name:          originalName,
		}
		var ret RestType
		if isWritableWithUnwritableSuperClass(classType, ownType) {
			uuidType := &restFindByUUIDClassType{
				superClass: superType,
				name:       originalName,
				nestedType: classType,
			}
			uuidType.uuidProperty = &RestProperty{
				Parent:   uuidType,
				Name:     "uuid",
				Required: true,
				Type:     NewFindBaseByUUIDObjectType(uuidType),
			}
			ret = uuidType
		} else if _, ok := writableSubclassCounts[originalName]; ok {
			polymorphicBase := &restPolymorphicBaseClassType{
				nestedType: classType,
				subtypes:   make([]RestType, 0),
			}
			ret = polymorphicBase
		} else {
			ret = classType
		}

		if parentResourceInfo != nil {
			ret = &restSubresourceClassType{
				name:       name,
				prefix:     parentResourceInfo.prefix,
				nestedType: ret,
			}
		}

		types[name] = ret
		classType.properties = buildProperties(classType, originalName, ownType, types)

		if polymorphicBaseType != nil {
			polyType := types[*polymorphicBaseType].(*restPolymorphicBaseClassType)
			polyType.subtypes = append(polyType.subtypes, classType)
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

func buildProperties(parent *restClassType, baseTypeName string, schema *openapi3.SchemaRef, types map[string]RestType) []*RestProperty {
	required := schema.Value.Required
	ret := make([]*RestProperty, 0)
	for name, property := range schema.Value.Properties {
		if name == "$type" {
			continue
		}
		if name == "additionalObjects" && baseTypeName == "authInternalAccount" {
			continue
		}
		rsSchemaTemplateBase := buildRSSchemaTemplateBase(schema, baseTypeName, name)
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
		restProperty.Type = buildType(baseTypeName, name, property, types, restProperty, rsSchemaTemplateBase)
		ret = append(ret, restProperty)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})
	return ret
}

func buildType(baseTypeName string, propertyName string, ref *openapi3.SchemaRef, types map[string]RestType, restProperty *RestProperty, rsSchemaTemplateBase map[string]any) RestPropertyType {
	schema := ref.Value
	if len(schema.AllOf) > 0 {
		if ref.Ref == "" {
			ref = schema.AllOf[0]
		}
		schema = schema.AllOf[0].Value
	}
	if schema.Type == "array" {
		return NewRestArrayType(buildType(baseTypeName, propertyName, schema.Items, types, restProperty, rsSchemaTemplateBase), rsSchemaTemplateBase)
	}
	if ref.Ref != "" && schema.Type == "string" && len(schema.Enum) > 0 {
		enumName := refToName(ref.Ref)
		enum := getOrBuildTypeModel(types, enumName, ref, nil)
		return NewEnumPropertyType(enum, rsSchemaTemplateBase)
	}
	if schema.Type == "boolean" || schema.Type == "integer" || schema.Type == "string" {
		return NewRestSimpleType(restProperty, schema.Type, schema.Format, rsSchemaTemplateBase)
	}
	if is(ref, object) {
		nestedTypeName := refToName(ref.Ref)
		if nestedTypeName == "" {
			nestedTypeName = baseTypeName + "_" + propertyName
		}
		nested := getOrBuildTypeModel(types, nestedTypeName, ref, nil)
		ret := NewNestedObjectType(restProperty, nested, rsSchemaTemplateBase)
		if ref.Ref != "" && is(ref, withUUID) && nested.Extends("Linkable") {
			if strings.HasSuffix(nestedTypeName, "Primer") {
				ret = NewFindByUUIDObjectType(ret, rsSchemaTemplateBase)
			}
		}
		return ret
	}

	log.Fatalf("Cannot construct a type for (%v+)", ref)
	return nil
}

func is(ref *openapi3.SchemaRef, check func(*openapi3.Schema) bool) bool {
	if check(ref.Value) {
		return true
	}
	for _, part := range ref.Value.AllOf {
		if check(part.Value) {
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
	return schema.Type == "object"
}

func deprecated(schema *openapi3.Schema) bool {
	return schema.Deprecated
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
func buildRSSchemaTemplateBase(ref *openapi3.SchemaRef, typeName string, propertyName string) map[string]any {
	required := slices.Contains(ref.Value.Required, propertyName)
	property := ref.Value.Properties[propertyName]
	if property.Value.AllOf != nil && len(property.Value.AllOf) > 0 {
		property = property.Value.AllOf[1]
	}
	defaultVal := property.Value.Extensions["x-tkh-default"]
	readOnly := property.Value.ReadOnly
	immutable := property.Value.Extensions["x-tkh-immutable"] != nil && property.Value.Extensions["x-tkh-immutable"].(bool)
	createOnly := property.Value.Extensions["x-tkh-create-only"] != nil && property.Value.Extensions["x-tkh-create-only"].(bool)

	if typeName == "Linkable" && (propertyName == "links" || propertyName == "permissions") {
		immutable = true
	}
	if immutable {
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
