package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println(colorErr("No argument was passed in."))
		printPossibleArgs()
		return
	}

	cmd := os.Args[1]
	f, ok := args[cmd]
	if !ok {
		fmt.Println(colorErr("An invalid argument was passed in."))
		printPossibleArgs()
		return
	}

	err := f(os.Args[2:])
	if err != nil {
		fmt.Println(colorErr("An error occurred:"))
		fmt.Println("  ", err)
	}

}

var args = map[string]func(args []string) error{
	"start":          Start,
	"setup-crdb":     SetupCockroach,
	"start-crdb":     StartWaitCRDB,
	"srcs":           ShowAdminSources,
	"deploy":         Deploy,
	"gen-env":        GenerateEnvFiles,
	"gen-admin-vers": GenAdminVersions,
}

func printPossibleArgs() {
	fmt.Println("The available arguments:")
	for k := range args {
		fmt.Println(" ", k)
	}
}
