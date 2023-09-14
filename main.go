package main

import (
	"bytes"
	"context"
	"embed"
	"go/format"
	"log"
	"os"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	keyhubsdk "github.com/topicuskeyhub/sdk-go"
	apimodel "github.com/topicuskeyhub/terraform-provider-keyhub-generator/model"
)

//go:embed templates/*
var tmpls embed.FS

func merge(file string, t *template.Template, model map[string]apimodel.RestType) {
	f, err := os.Create("./" + file + ".go")
	if err != nil {
		log.Fatalf("cannot create %s: %s", file, err)
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, file+".go.tmpl", model)
	if err != nil {
		log.Fatalf("Template %s failed: %s", file, err)
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
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: false}
	doc, err := loader.LoadFromData(keyhubsdk.OpenAPISpec())
	if err != nil {
		log.Fatalf("Cannot read openapi file: %s", err)
	}
	err = doc.Validate(ctx, openapi3.DisableExamplesValidation())
	if err != nil {
		log.Fatalf("Validation of openapi file failed: %s", err)
	}

	functions := template.FuncMap{
		"RecurseCutOff":             apimodel.RecurseCutOff,
		"AdditionalObjectsProperty": apimodel.AdditionalObjectsProperty,
		"AllDirectProperties":       apimodel.AllDirectProperties,
	}
	model := apimodel.BuildModel(doc)
	t, err := template.New("provider").Funcs(functions).ParseFS(tmpls, "templates/*")
	if err != nil {
		log.Fatalf("Template parsing failed: %s", err)
	}

	merge("full-data-struct-ds", t, model)
	merge("full-data-struct-rs", t, model)
	merge("full-helpers", t, model)
	merge("full-object-attrs-ds", t, model)
	merge("full-object-attrs-rs", t, model)
	merge("full-schema-ds", t, model)
	merge("full-schema-rs", t, model)
	merge("full-tf-to-data-struct-ds", t, model)
	merge("full-tf-to-data-struct-rs", t, model)
	merge("full-tf-to-tkh-ds", t, model)
	merge("full-tf-to-tkh-rs", t, model)
	merge("full-tkh-to-tf-ds", t, model)
	merge("full-tkh-to-tf-rs", t, model)
}
