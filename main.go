package main

import (
    "github.com/go-openapi/loads"
    "github.com/go-openapi/spec"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "flag"
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

    b, err := json.MarshalIndent(primary, "", "  ")
    check(err)
    out := *outPtr;
    if out == "" {
        out = *inPtr
    }

    we := ioutil.WriteFile(out, b, 0644)
    check(we)
}

func GenSwagger() *spec.Swagger {
    swagger := &spec.Swagger{
        SwaggerProps: spec.SwaggerProps{
            SecurityDefinitions: map[string] *spec.SecurityScheme{
                "ApiKey": spec.APIKeyAuth("api_key", "header"),
            },
            Security: []map[string][]string{
                {"ApiKey": {}},
            },
        },
    }
    return swagger
}
