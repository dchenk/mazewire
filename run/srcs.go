package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ShowAdminSources(_ []string) error {
	const dir = "front/dist"
	//fileVersions := make([]string, 0, 4)
	//err := filepath.Walk(dir, visitStaticDir(&files))
	//if err != nil {
	//	return err
	//}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	fmt.Println("{")
	for i, file := range files {
		if i > 0 {
			fmt.Print(",\n")
		}
		split := strings.Split(file.Name(), ".")
		fmt.Printf("  %q: %q", "admin-"+split[0], split[1])
	}
	fmt.Println("\n}")
	return nil
}

func visitStaticDir(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		*files = append(*files, path)
		return nil
	}
}
