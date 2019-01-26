// +build !test

// Note: This file is not built with the rest of the package for unit tests.
// To test this package alone, do: go test -v -tags=test

package env

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.Parse()
	var err error
	switch *flagEnv {
	case envConfigFile:
		vars, err = parseFile(*flagPath)
	case envConfigVars:
		vars, err = parseVars()
	default:
		fmt.Printf("The \"env\" argument must be either %q or %q; got unexpected value %q\n", envConfigFile, envConfigVars, *flagEnv)
		os.Exit(2)
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
