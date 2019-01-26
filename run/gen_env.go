package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

const (
	envProdGoFile = "main/env.go"
	envDevGoFile  = "main/env_dev.go"
)

func GenerateEnvFiles(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("specify if to generate the environment variables file for 'prod' or 'dev'")
	}

	var buf bytes.Buffer
	buf.WriteString("// +build ")

	var varsFile, outFile string

	switch env := args[0]; env {
	case "prod":
		buf.WriteString("prod\n\n")
		varsFile = envVarsProdFile
		outFile = envProdGoFile
	case "dev":
		buf.WriteString("!prod\n\n")
		varsFile = envVarsDevFile
		outFile = envDevGoFile
	default:
		return fmt.Errorf("unexpected environment %q", env)
	}

	envBytes, err := ioutil.ReadFile(varsFile)
	if err != nil {
		return fmt.Errorf("could not read env file: %v", err)
	}

	lines := bytes.Split(envBytes, []byte{'\n'})

	buf.WriteString("package main\n\nconst (\n")

	for len(lines) > 0 {
		if len(lines[0]) == 0 {
			// Skip the blank line
			lines = lines[1:]
			break
		}
		buf.WriteByte('\t')
		buf.Write(lines[0])
		buf.WriteByte('\n')
		lines = lines[1:]
	}

	buf.WriteString(")\n\nvar (\n")

	for len(lines) > 0 {
		buf.WriteByte('\t')
		buf.Write(lines[0])
		buf.WriteByte('\n')
		lines = lines[1:]
	}

	buf.WriteString(")\n") // gofmt adds a trailing newline

	return ioutil.WriteFile(outFile, buf.Bytes(), 644)
}
