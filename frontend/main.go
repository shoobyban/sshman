package main

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed dist
var files embed.FS

func main() {
	templates, _ := fs.ReadDir(files, "dist")
	for _, template := range templates {
		fmt.Printf("%q\n", template.Name())
	}
}
