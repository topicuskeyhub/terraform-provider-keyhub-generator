package main

import (
	"context"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
	keyhubsdk "github.com/topicuskeyhub/sdk-go"
	apimodel "github.com/topicuskeyhub/terraform-provider-keyhub-generator/model"
)

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
	tmp := model["directoryAccountDirectory"]
	log.Printf(("(%+v)"), tmp)
}
