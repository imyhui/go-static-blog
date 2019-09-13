package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

type Post struct {
	Title   string
	Date    string
	Content string
	Source  []byte
	URL     string
}

func getSources() []string {
	files, _ := filepath.Glob("srcs/*.md")
	return files
}

func renderMarkdown(source []byte) string {
	content := string(blackfriday.Run(source))
	return content
}

func parseSource(fileName string) Post {
	sources, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(sources), "\n")
	title := strings.Split(string(lines[1]), ": ")[1]
	date := strings.Split(string(lines[2]), ": ")[1]
	source := []byte(strings.Join(lines[5:len(lines)], "\n"))
	content := renderMarkdown(source)
	URL := strings.Replace(strings.ToLower(title), " ", "-", -1)
	return Post{title, date, content, source, URL}
}

func main() {
	files := getSources()
	for _, file := range files {
		post := parseSource(file)
		fmt.Printf("%v\n", post)
	}
}
