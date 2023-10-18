package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"go/format"
	"log"
	"os"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	keyhubsdk "github.com/topicuskeyhub/sdk-go"
	apimodel "github.com/topicuskeyhub/terraform-provider-keyhub-generator/model"
)

var resource = flag.String("resource", "", "Generate a data source or resource for the given SDK resource, ie. client")
var linkable = flag.String("linkable", "", "The full name of the linkable type, ie. clientClientApplication")
var mode = flag.String("mode", "", "'model', 'data' or 'resource'")

//go:embed templates/*
var tmpls embed.FS

type resourceTemplateParameters struct {
	UpdateSupported                    bool
	DeleteSupported                    bool
	ResourceBase                       string
	ResourceBaseUp                     string
	Name                               string
	NameUp                             string
	NameUnderscore                     string
	FullName                           string
	FullNameUp                         string
	BaseName                           string
	BaseNameUp                         string
	CollectionRequestTypePrefix        string
	ItemRequestTypePrefix              string
	SubResourceReqMethod               string
	SubResourceBaseUp                  string
	ParentResourceType                 string
	ParentResourceNamePrefixUp         string
	ParentResourceNamePrefixUnderscore string
}

var resourceTemplateConfigs = map[string]resourceTemplateParameters{
	"clientapplication": {
		UpdateSupported:                    false,
		DeleteSupported:                    false,
		ResourceBase:                       "client",
		ResourceBaseUp:                     "Client",
		Name:                               "clientapplication",
		NameUp:                             "Clientapplication",
		NameUnderscore:                     "clientapplication",
		FullName:                           "clientClientApplication",
		FullNameUp:                         "ClientClientApplication",
		BaseName:                           "clientClientApplication",
		BaseNameUp:                         "ClientClientApplication",
		CollectionRequestTypePrefix:        "Client",
		ItemRequestTypePrefix:              "WithClientItem",
		SubResourceReqMethod:               "",
		SubResourceBaseUp:                  "",
		ParentResourceType:                 "",
		ParentResourceNamePrefixUp:         "",
		ParentResourceNamePrefixUnderscore: "",
	},
	"client_vaultrecord": {
		UpdateSupported:                    true,
		DeleteSupported:                    true,
		ResourceBase:                       "client",
		ResourceBaseUp:                     "Client",
		Name:                               "clientVaultrecord",
		NameUp:                             "ClientVaultrecord",
		NameUnderscore:                     "client_vaultrecord",
		FullName:                           "clientApplicationVaultVaultRecord",
		FullNameUp:                         "ClientApplicationVaultVaultRecord",
		BaseName:                           "vaultVaultRecord",
		BaseNameUp:                         "VaultVaultRecord",
		CollectionRequestTypePrefix:        "ItemVaultRecord",
		ItemRequestTypePrefix:              "ItemVaultRecordWithRecordItem",
		SubResourceReqMethod:               ".Vault().Record()",
		SubResourceBaseUp:                  "Record",
		ParentResourceType:                 "ClientClientApplicationPrimer",
		ParentResourceNamePrefixUp:         "ClientApplication",
		ParentResourceNamePrefixUnderscore: "client_application",
	},
	"group_vaultrecord": {
		UpdateSupported:                    true,
		DeleteSupported:                    true,
		ResourceBase:                       "group",
		ResourceBaseUp:                     "Group",
		Name:                               "groupVaultrecord",
		NameUp:                             "GroupVaultrecord",
		NameUnderscore:                     "group_vaultrecord",
		FullName:                           "groupVaultVaultRecord",
		FullNameUp:                         "GroupVaultVaultRecord",
		BaseName:                           "vaultVaultRecord",
		BaseNameUp:                         "VaultVaultRecord",
		CollectionRequestTypePrefix:        "ItemVaultRecord",
		ItemRequestTypePrefix:              "ItemVaultRecordWithRecordItem",
		SubResourceReqMethod:               ".Vault().Record()",
		SubResourceBaseUp:                  "Record",
		ParentResourceType:                 "GroupGroupPrimer",
		ParentResourceNamePrefixUp:         "Group",
		ParentResourceNamePrefixUnderscore: "group",
	},
	"group": {
		UpdateSupported:                    false,
		DeleteSupported:                    false,
		ResourceBase:                       "group",
		ResourceBaseUp:                     "Group",
		Name:                               "group",
		NameUp:                             "Group",
		NameUnderscore:                     "group",
		FullName:                           "groupGroup",
		FullNameUp:                         "GroupGroup",
		BaseName:                           "groupGroup",
		BaseNameUp:                         "GroupGroup",
		CollectionRequestTypePrefix:        "Group",
		ItemRequestTypePrefix:              "WithGroupItem",
		SubResourceReqMethod:               "",
		SubResourceBaseUp:                  "",
		ParentResourceType:                 "",
		ParentResourceNamePrefixUp:         "",
		ParentResourceNamePrefixUnderscore: "",
	},
	"grouponsystem": {
		UpdateSupported:                    false,
		DeleteSupported:                    false,
		ResourceBase:                       "system",
		ResourceBaseUp:                     "System",
		Name:                               "grouponsystem",
		NameUp:                             "Grouponsystem",
		NameUnderscore:                     "grouponsystem",
		FullName:                           "nestedProvisioningGroupOnSystem",
		FullNameUp:                         "NestedProvisioningGroupOnSystem",
		BaseName:                           "provisioningGroupOnSystem",
		BaseNameUp:                         "ProvisioningGroupOnSystem",
		CollectionRequestTypePrefix:        "ItemGroup",
		ItemRequestTypePrefix:              "ItemGroupWithGroupItem",
		SubResourceReqMethod:               ".Group()",
		SubResourceBaseUp:                  "Group",
		ParentResourceType:                 "ProvisioningProvisionedSystemPrimer",
		ParentResourceNamePrefixUp:         "ProvisionedSystem",
		ParentResourceNamePrefixUnderscore: "provisioned_system",
	},
	"serviceaccount": {
		UpdateSupported:                    true,
		DeleteSupported:                    false,
		ResourceBase:                       "serviceaccount",
		ResourceBaseUp:                     "Serviceaccount",
		Name:                               "serviceaccount",
		NameUp:                             "Serviceaccount",
		NameUnderscore:                     "serviceaccount",
		FullName:                           "serviceaccountServiceAccount",
		FullNameUp:                         "ServiceaccountServiceAccount",
		BaseName:                           "serviceaccountServiceAccount",
		BaseNameUp:                         "ServiceaccountServiceAccount",
		CollectionRequestTypePrefix:        "Serviceaccount",
		ItemRequestTypePrefix:              "WithServiceaccountItem",
		SubResourceReqMethod:               "",
		SubResourceBaseUp:                  "",
		ParentResourceType:                 "",
		ParentResourceNamePrefixUp:         "",
		ParentResourceNamePrefixUnderscore: "",
	},
}

func merge(template string, suffix string, t *template.Template, model any) {
	file := template + suffix
	f, err := os.Create("internal/provider/" + file + ".go")
	if err != nil {
		log.Fatalf("cannot create %s: %s", file, err)
	}
	log.Printf(" ... writing %s", f.Name())

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, template+".go.tmpl", model)
	if err != nil {
		log.Fatalf("Template %s failed: %s", template, err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("Format %s failed: %s", file, err)
		p = buf.Bytes()
	}
	f.Write(p)
	f.Close()
}

func main() {
	flag.Parse()
	log.Println("Generating Topicus KeyHub Terraform Provider source...")
	ctx := context.Background()

	if *mode == "model" {
		functions := template.FuncMap{
			"RecurseCutOff":             apimodel.RecurseCutOff,
			"AdditionalObjectsProperty": apimodel.AdditionalObjectsProperty,
			"AllDirectProperties":       apimodel.AllDirectProperties,
		}
		t, err := template.New("provider").Funcs(functions).ParseFS(tmpls, "templates/model/*")
		if err != nil {
			log.Fatalf("Template parsing failed: %s", err)
		}

		loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: false}
		doc, err := loader.LoadFromData(keyhubsdk.OpenAPISpec())
		if err != nil {
			log.Fatalf("Cannot read openapi file: %s", err)
		}
		err = doc.Validate(ctx, openapi3.DisableExamplesValidation())
		if err != nil {
			log.Fatalf("Validation of openapi file failed: %s", err)
		}
		model := apimodel.BuildModel(doc)
		merge("full-data-struct-ds", "", t, model)
		merge("full-data-struct-rs", "", t, model)
		merge("full-helpers", "", t, model)
		merge("full-object-attrs-ds", "", t, model)
		merge("full-object-attrs-rs", "", t, model)
		merge("full-schema-ds", "", t, model)
		merge("full-schema-rs", "", t, model)
		merge("full-tf-to-data-struct-ds", "", t, model)
		merge("full-tf-to-data-struct-rs", "", t, model)
		merge("full-tf-to-tkh-ds", "", t, model)
		merge("full-tf-to-tkh-rs", "", t, model)
		merge("full-tkh-to-tf-ds", "", t, model)
		merge("full-tkh-to-tf-rs", "", t, model)
	} else if *mode == "data" {
		t, err := template.New("provider").ParseFS(tmpls, "templates/impl/*")
		if err != nil {
			log.Fatalf("Template parsing failed: %s", err)
		}
		merge("datasource", "-"+*resource, t, map[string]string{
			"Name":           *resource,
			"NameUp":         apimodel.FirstCharToUpper(*resource),
			"FullName":       *linkable,
			"FullNameUp":     apimodel.FirstCharToUpper(*linkable),
			"ResourceBase":   *resource,
			"ResourceBaseUp": apimodel.FirstCharToUpper(*resource),
		})
	} else if *mode == "resource" {
		t, err := template.New("provider").ParseFS(tmpls, "templates/impl/*")
		if err != nil {
			log.Fatalf("Template parsing failed: %s", err)
		}
		merge("resource", "-"+*resource, t, resourceTemplateConfigs[*resource])
	}
}
