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

	model := apimodel.BuildModel(doc)
	t, err := template.New("provider").ParseFS(tmpls, "templates/*")
	if err != nil {
		log.Fatalf("Template parsing failed: %s", err)
	}

	f, err := os.Create("./full.go")
	if err != nil {
		log.Fatalf("cannot create file: %s", err)
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "full.go.tmpl", model)
	if err != nil {
		log.Fatalf("Template failed: %s", err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("Format failed: %s", err)
		p = buf.Bytes()
	}
	f.Write(p)

	f.Close()
}
