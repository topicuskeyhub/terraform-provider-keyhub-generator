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

	if *mode == "impl" {
		functions := template.FuncMap{
			"RecurseCutOff":             apimodel.RecurseCutOff,
			"AdditionalObjectsProperty": apimodel.AdditionalObjectsProperty,
			"AllDirectProperties":       apimodel.AllDirectProperties,
		}
		t, err := template.New("provider").Funcs(functions).ParseFS(tmpls, "templates/impl/*")
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
		t, err := template.New("provider").ParseFS(tmpls, "templates/data/*")
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
	}
}
