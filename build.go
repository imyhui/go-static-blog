package main

import (
	"fmt"
	"path/filepath"
)

func getSources() []string {
	files, _ := filepath.Glob("srcs/*.md")
	return files
}

func main() {
	files := getSources()
	fmt.Printf("%v", files)
}
