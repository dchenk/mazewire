package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
)

func prefixLines(r io.Reader, prefix string) []byte {
	liner := bufio.NewReader(r)
	buffer := bytes.Buffer{}
	for {
		line, err := liner.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		buffer.WriteString(prefix)
		buffer.Write(line)
	}
	return buffer.Bytes()
}

const componentsSchemas = "#/components/schemas"

func GenerateOpenAPI(_ []string) error {
	tables, err := readTables()
	if err != nil {
		return err
	}

	tablesSchemas, err := yaml.Marshal(tables)
	if err != nil {
		return fmt.Errorf("could not marshal tables schemas; %v", err)
	}

	const indent = "    "

	file, err := os.Open("api/openapi.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.Buffer{}

	didTables := false
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("could not read line in tables file; %v", err)
		}
		if !didTables && bytes.Contains(line, []byte("# Generated schema code")) {
			buffer.Write(prefixLines(bytes.NewBuffer(tablesSchemas), indent))
			didTables = true
		} else {
			buffer.Write(line)
		}
	}

	loader := openapi3.SwaggerLoader{}

	spec, err := loader.LoadSwaggerFromData(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not load spec file; %v", err)
	}

	if err := spec.Validate(context.Background()); err != nil {
		return fmt.Errorf("spec is not valid: %v", err)
	}

	specJSON, err := spec.MarshalJSON()
	if err != nil {
		return fmt.Errorf("could not marshal JSON; %v", err)
	}

	fmt.Printf("%s\n", specJSON)
	return nil
}

func readTables() (TypesDoc, error) {
	file, err := os.Open("db/tables.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not open tables file; %v", err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)
	data := make(TypesDoc, 10)
	return data, decoder.Decode(&data)
}

// integerType matches integer type names that have an explicit size.
var integerType = regexp.MustCompile("^int\\d+$")

// The TypesDoc type is the type of the "tables.yaml" file.
type TypesDoc map[string]Schema

type Schema struct {
	Type           string              `yaml:"type"`
	Format         string              `yaml:"format,omitempty"`
	Description    string              `yaml:"description"`
	Properties     map[string]Property `yaml:"properties,omitempty"`
	Enum           []EnumOption        `yaml:"enum,omitempty"`
	ProtoPackage   string              `yaml:"x-proto-package,omitempty"`
	ProtoGoPackage string              `yaml:"x-proto-go_package,omitempty"`
}

// MarshalYAML implements yaml.Marshaler to return an object that can be a valid OpenAPI schema.
func (s Schema) MarshalYAML() (interface{}, error) {
	if integerType.MatchString(s.Type) {
		s.Format = s.Type
		s.Type = "integer"
	}
	s.ProtoGoPackage = ""
	return s, nil
}

// A Property is a property of a Schema.
type Property struct {
	// Ref can point to a Schema, and if it's set then no other fields should be set.
	Ref         *Ref         `yaml:"$ref,omitempty"`
	Type        string       `yaml:"type,omitempty"`
	Format      string       `yaml:"format,omitempty"`
	Description string       `yaml:"description,omitempty"`
	Enum        []EnumOption `yaml:"enum,omitempty"`
	ProtoField  uint32       `yaml:"x-proto-field,omitempty"`
}

// MarshalYAML implements yaml.Marshaler to return an object that can be a valid OpenAPI property.
func (p Property) MarshalYAML() (interface{}, error) {
	if integerType.MatchString(p.Type) {
		p.Format = p.Type
		p.Type = "integer"
	}
	p.ProtoField = 0
	return p, nil
}

// A Ref locates another schema.
type Ref string

// MarshalYAML for the Ref type expects the ref string to start with "#/" and point to a schema
// at the root of a TypesDoc file.
func (r Ref) MarshalYAML() (interface{}, error) {
	return componentsSchemas + strings.TrimPrefix(string(r), "#"), nil
}

type EnumOption struct {
	Value     uint32 `yaml:"value"`
	ProtoName string `yaml:"x-proto-name"`
}

func (eo EnumOption) MarshalYAML() (interface{}, error) {
	return eo.Value, nil
}
