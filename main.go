package main

import (
	"context"
	"embed"
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

	for _, curType := range model["groupGroup"].AllRequiredTypes() {
		err = t.ExecuteTemplate(os.Stdout, "datastruct.go.tmpl", curType)
		if err != nil {
			log.Fatalf("Template failed: %s", err)
		}
	}
}
