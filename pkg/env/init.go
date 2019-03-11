// +build !test

// Note: This file is not built with the rest of the package for unit tests.
// To test this package alone, do: go test -v -tags=test

package env

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// init parses the environment variables and verifies that all the required variables are set.
// If there is an error, a message is printed and the program exits with status 2.
func init() {
	flag.Parse()

	*flagEnvFile = strings.TrimSpace(*flagEnvFile)

	var err error
	if *flagEnvFile == "" {
		vars, err = parseVars()
	} else {
		vars, err = parseFile(*flagEnvFile)
	}

	if err != nil {
		fmt.Printf("An error occurred determining the environment variables; %v\n", err)
		os.Exit(2)
	}

	if missing := checkForRequiredVars(vars); len(missing) > 0 {
		if len(missing) == 1 {
			fmt.Printf("The required environment variable %v is missing\n", missing[0])
		} else {
			fmt.Println("The following required environment variables are missing:")
			for _, miss := range missing {
				fmt.Printf("\t- %v\n", miss)
			}
		}
		os.Exit(2)
	}
}
