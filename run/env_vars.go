package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

// The CWD when Make is invoked is one directory higher.
const (
	envVarsDevFile  = ".sample.env"
	envVarsProdFile = ".prod.env"
)

var envVars = make(map[string]string, 8)

func init() {

	envBytes, err := ioutil.ReadFile(envVarsDevFile)
	if err != nil {
		fmt.Printf("Could not read env file: %v\n", err)
		return
	}

	lines := bytes.Split(envBytes, []byte{'\n'})

	key := make([]byte, 0, 24)
	for _, v := range lines {
		v = bytes.TrimSpace(v)
		if len(v) == 0 {
			// When we get to the blank line separating the actual constants from the variables,
			// skip the remaining variables.
			break
		}
		key = key[:0]
		// Extract the name of the variable.
		for i := 0; i < len(v); i++ {
			char := v[i]
			if (char >= 'A' && char <= 'Z') || char == '_' {
				key = append(key, char)
			} else {
				break
			}
		}
		vk := string(key)
		v = bytes.TrimLeft(v, vk+" =\t")           // First just the beginning
		envVars[vk] = strings.Trim(string(v), `"`) // Now the quotes around the connection URL
	}

}

// CockroachPort returns the port that CockroachDB is supposed to run on in dev mode.
// The port is the string following the last ':' in DB_CONNECTION.
func CockroachPort() string {
	return envVars["DB_CONNECTION"][strings.LastIndexByte(envVars["DB_CONNECTION"], ':')+1:]
}
