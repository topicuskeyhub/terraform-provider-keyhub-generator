// Copyright (c) Topicus Security B.V.
// SPDX-License-Identifier: APSL-2.0

package model

import (
	"strings"
	"unicode"
)

type RestProperty struct {
	Parent     RestType
	Name       string
	Type       RestPropertyType
	Required   bool
	WriteOnly  bool
	Deprecated bool
	dsProperty *RestProperty
}

type RestPropertyType interface {
	MarkReachable()
	PropertyNameSuffix() string
	FlattenMode() string
	TFName() string
	TFValueType() string
	TFAttrType(inAdditionalObjects bool) string
	TFValidatorType() string
	TFValidators() []string
	ToTFAttrWithDiag() bool
	ToTKHAttrWithDiag() bool
	ToTKHCustomCode() string
	TFAttrNeeded() bool
	Complex() bool
	NestedType() RestType
	TKHToTF(value string, listItem bool) string
	TFToTKH(value string, listItem bool) string
	TKHToTFGuard() string
	TFToTKHGuard() string
	TKHGetter(propertyName string) string
	SDKTypeName(listItem bool) string
	SDKTypeConstructor() string
	DSSchemaTemplate() string
	DSSchemaTemplateData() map[string]any
	RSSchemaTemplate() string
	RSSchemaTemplateData() map[string]any
	DS() RestPropertyType
}

var aliasses = map[string]string{
	"auditAuditRecord.auditAuditRecordType":                                       "type",
	"certificateCertificatePrimer.certificateCertificatePrimerType":               "type",
	"clientClientApplicationPrimer.clientClientApplicationPrimerType":             "type",
	"directoryAccountDirectoryPrimer.directoryAccountDirectoryPrimerType":         "type",
	"directoryAccountDirectorySummary.directoryAccountDirectorySummaryType":       "type",
	"provisioningGroupOnSystemPrimer.provisioningGroupOnSystemPrimerType":         "type",
	"provisioningProvisionedSystemPrimer.provisioningProvisionedSystemPrimerType": "type",
	"markItemMarker.markItemMarkerType":                                           "type",
	"vaultVaultRecordShare.vaultVaultRecordShareType":                             "type",
	"webhookWebhookPush.webhookWebhookPushType":                                   "type",
}

func (p *RestProperty) internalName() string {
	return p.Name + p.Type.PropertyNameSuffix()
}

func (p *RestProperty) GoName() string {
	ret := FirstCharToUpper(p.internalName())
	ret = strings.ReplaceAll(ret, "Oidc", "OIDC")
	ret = strings.ReplaceAll(ret, "Oauth2", "OAuth2")
	ret = strings.ReplaceAll(ret, "Ldap", "LDAP")
	ret = strings.ReplaceAll(ret, "Scim", "SCIM")
	ret = strings.ReplaceAll(ret, "Uuid", "UUID")
	ret = strings.ReplaceAll(ret, "Uid", "UID")
	ret = strings.ReplaceAll(ret, "Id", "ID")
	ret = strings.ReplaceAll(ret, "Rdn", "RDN")
	ret = strings.ReplaceAll(ret, "Dn", "DN")
	ret = strings.ReplaceAll(ret, "Url", "URL")
	ret = strings.ReplaceAll(ret, "Uri", "URI")
	ret = strings.ReplaceAll(ret, "Tls", "TLS")
	return ret
}

func (p *RestProperty) TFName() string {
	name := p.internalName()
	if alias, ok := aliasses[p.Parent.APITypeName()+"."+name]; ok {
		name = alias
	}
	name = strings.ReplaceAll(name, "OIDC", "Oidc")
	name = strings.ReplaceAll(name, "OAuth2", "Oauth2")
	name = strings.ReplaceAll(name, "LDAP", "Ldap")
	name = strings.ReplaceAll(name, "SCIM", "Scim")
	name = strings.ReplaceAll(name, "UUID", "Uuid")
	name = strings.ReplaceAll(name, "UID", "Uid")
	name = strings.ReplaceAll(name, "ID", "Id")
	name = strings.ReplaceAll(name, "RDN", "Rdn")
	name = strings.ReplaceAll(name, "DN", "Dn")
	name = strings.ReplaceAll(name, "URL", "Url")
	name = strings.ReplaceAll(name, "URI", "Uri")
	name = strings.ReplaceAll(name, "TLS", "Tls")
	// modify one way only
	name = strings.ReplaceAll(name, "2FA", "2fa")
	name = strings.ReplaceAll(name, "FA", "Fa")
	ret := make([]rune, 0)
	for _, r := range name {
		if unicode.IsUpper(r) {
			ret = append(ret, '_', unicode.ToLower(r))
		} else {
			ret = append(ret, r)
		}
	}
	return string(ret)
}

func (p *RestProperty) TFType() string {
	return p.Type.TFName()
}

func (p *RestProperty) TFAttrType(inAdditionalObjects bool) string {
	return p.Type.TFAttrType(inAdditionalObjects)
}

func (p *RestProperty) TKHToTF() string {
	return p.Type.TKHToTF(p.Type.TKHGetter(p.Name), false)
}

func (p *RestProperty) TKHSetter() string {
	return "Set" + FirstCharToUpper(p.Name)
}

func (p *RestProperty) TFToTKH() string {
	return p.Type.TFToTKH("objAttrs[\""+p.TFName()+"\"]", false)
}

func (p *RestProperty) DS() *RestProperty {
	if p.dsProperty != nil {
		// break recursion
		return p.dsProperty
	}

	p.dsProperty = &RestProperty{
		Parent:     p.Parent,
		Name:       p.Name,
		Required:   p.Required,
		WriteOnly:  p.WriteOnly,
		Deprecated: p.Deprecated,
	}
	p.dsProperty.Type = p.Type.DS()
	return p.dsProperty
}

func (p *RestProperty) IsDTypeRequired() bool {
	return strings.HasSuffix(p.Parent.GoTypeName(), "_additionalObjects") &&
		p.Type.NestedType() != nil && p.Type.NestedType().APIDiscriminator() != "" &&
		!p.Type.NestedType().Extends("Linkable") && !p.Type.NestedType().Extends("NonLinkable")
}
