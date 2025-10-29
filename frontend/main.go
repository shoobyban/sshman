package main

import (
	"embed"
	"fmt"
	"io/fs"
)

// Note: embedding dist at compile time requires the dist directory to exist.
// Leave files empty when dist is not present so tests and builds don't fail.
var files embed.FS

func main() {
	templates, err := fs.ReadDir(files, "dist")
	if err != nil {
		// no embedded assets available, exit gracefully
		fmt.Println("no embedded assets")
		return
	}
	for _, template := range templates {
		fmt.Printf("%q\n", template.Name())
	}
}
