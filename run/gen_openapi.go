package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
)

func GenerateOpenAPI(_ []string) error {
	return readTables()
	data, err := ioutil.ReadFile("openapi.yaml")
	if err != nil {
		return fmt.Errorf("could not read spec file; %v", err)
	}

	loader := openapi3.SwaggerLoader{
		IsExternalRefsAllowed: true,
		// LoadSwaggerFromURIFunc:
		// func(loader *openapi3.SwaggerLoader, url *url.URL) (*openapi3.Swagger, error) {
		// 	fmt.Print(loader, url)
		// 	return nil, nil
		// },
	}
	spec, err := loader.LoadSwaggerFromData(data)
	if err != nil {
		return fmt.Errorf("could not load spec file; %v", err)
	}
	fmt.Print(spec)
	return nil
}

func readTables() error {
	file, err := os.Open("db/tables.yaml")
	if err != nil {
		return fmt.Errorf("could not open tables file; %v", err)
	}
	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(true)
	data := make(TablesDoc, 4)
	err = decoder.Decode(&data)
	if err != nil {
		return fmt.Errorf("DECODING; %v", err)
	}
	fmt.Printf("%#v\n", data)
	return nil
}

type OpenAPISpecElement interface {
	// ToOAS returns the YAML that specifies the permissible values for the element, which is the
	// type schema for Schemas or Properties or the literal value for EnumOptions.
	// The indent string must be placed in front of each line.
	ToOAS(indent string) string
}

// The TablesDoc type is the type of the "tables.yaml" file.
type TablesDoc map[string]Schema

func (r TablesDoc) ToOAS(indent string) string {
	builder := strings.Builder{}
	for k, v := range r {
		builder.WriteString(indent)
		builder.WriteString(k)
		builder.WriteString(": ")
		builder.WriteString(v.ToOAS(indent))
	}
	return builder.String()
}

type Schema struct {
	Type        string              `yaml:"type"`
	Format      string              `yaml:"format"`
	Description string              `yaml:"description"`
	Properties  map[string]Property `yaml:"properties"`
	Enum        []EnumOption        `yaml:"enum"`
}

func (s *Schema) ToOAS(indent string) string {
	builder := strings.Builder{}
	err := yaml.NewEncoder(&builder).Encode(s)
	if err != nil {
		panic(fmt.Sprintf("could not encode schema; %v", err))
	}
	return builder.String()
}

type Property struct {
	Ref         *Ref   `yaml:"$ref"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Enum        Enum   `yaml:"enum"`
	ProtoField  uint32 `yaml:"x-proto-field"`
}

func (s *Property) ToOAS(indent string) string {
	builder := strings.Builder{}
	err := yaml.NewEncoder(&builder).Encode(s)
	if err != nil {
		panic(fmt.Sprintf("could not encode schema; %v", err))
	}
	return builder.String()
}

// A Ref locates another schema.
type Ref string

// ToOAS for the Ref type expects the ref string to start with "#/" and point to a schema within
// the same file.
func (r Ref) ToOAS(indent string) string {
	return indent + "#/components/schemas" + strings.TrimPrefix(string(r), "#")
}

type Enum []EnumOption

func (e Enum) ToOAS(indent string) string {
	builder := strings.Builder{}
	for _, v := range e {
		builder.WriteString(v.ToOAS(indent))
	}
	return builder.String()
}

type EnumOption struct {
	Value     uint32 `yaml:"value"`
	ProtoName string `yaml:"x-proto-name"`
}

func (eo *EnumOption) ToOAS(indent string) string {
	return fmt.Sprintf("- %s%d\n%s  # Field name: %s", indent, eo.Value, indent, eo.ProtoName)
}
