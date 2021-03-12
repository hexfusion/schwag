package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/go-openapi/loads"
)

func main() {
	var output string
	var input string

	// set flags
	flag.StringVar(&input, "input", "", "a string")
	flag.StringVar(&output, "output", "", "a string")
	flag.Parse()

	if output == "" {
		output = input
	}

	file, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}

	base, err := loads.Analyzed(json.RawMessage([]byte(file)), "")
	if err != nil {
		log.Fatal(err)
	}

	primary := base.Spec()
	// this is all we do
	primary.Schemes = []string{"http", "https"}

	b, err := json.MarshalIndent(primary, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(output, b, 0644); err != nil {
		log.Fatal(err)
	}
}
