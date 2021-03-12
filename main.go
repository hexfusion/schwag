package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"io/ioutil"
)

func check(err error) {
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}
}

func main() {
	// set flags
	inPtr := flag.String("input", "", "a string")
	outPtr := flag.String("output", "", "a string")
	flag.Parse()

	file, err := ioutil.ReadFile(*inPtr)
	check(err)

	base, err := loads.Analyzed(json.RawMessage([]byte(file)), "")
	check(err)

	primary := base.Spec()
	swagStub := GenSwagger()
	primary.Security = swagStub.Security
	primary.SecurityDefinitions = swagStub.SecurityDefinitions

	FixEmptyResponseDescriptions(primary)

	b, err := json.MarshalIndent(primary, "", "  ")
	check(err)
	out := *outPtr
	if out == "" {
		out = *inPtr
	}

	we := ioutil.WriteFile(out, b, 0644)
	check(we)
}

func GenSwagger() *spec.Swagger {
	swagger := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			SecurityDefinitions: map[string]*spec.SecurityScheme{
				"ApiKey": spec.APIKeyAuth("Authorization", "header"),
			},
			Security: []map[string][]string{
				{"ApiKey": {}},
			},
		},
	}
	return swagger
}

// thanks to github.com/msample/swagger-mixin for the empty desc solution
func FixEmptyResponseDescriptions(s *spec.Swagger) {
	for _, v := range s.Paths.Paths {
		if v.Put != nil {
			FixEmptyDescs(v.Put.Responses)
		}
		if v.Post != nil {
			FixEmptyDescs(v.Post.Responses)
		}
		if v.Delete != nil {
			FixEmptyDescs(v.Delete.Responses)
		}
		if v.Options != nil {
			FixEmptyDescs(v.Options.Responses)
		}
		if v.Head != nil {
			FixEmptyDescs(v.Head.Responses)
		}
		if v.Patch != nil {
			FixEmptyDescs(v.Patch.Responses)
		}
	}
	for k, v := range s.Responses {
		FixEmptyDesc(&v)
		s.Responses[k] = v
	}
}

// FixEmptyDescs adds "(empty)" as the description for any Response in
// the given Responses object that doesn't already have one.
func FixEmptyDescs(rs *spec.Responses) {
	FixEmptyDesc(rs.Default)
	for k, v := range rs.StatusCodeResponses {
		FixEmptyDesc(&v)
		rs.StatusCodeResponses[k] = v
	}
}

// FixEmptyDesc adds "A successful response." as the description to the given
// Response object if it doesn't already have one and isn't a
// ref. No-op on nil input.
func FixEmptyDesc(rs *spec.Response) {
	if rs == nil || rs.Description != "" || rs.Ref.Ref.GetURL() != nil {
		return
	}
	rs.Description = "A successful response."
}
