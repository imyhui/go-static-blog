package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/russross/blackfriday"
)

type Post struct {
	Title   string
	Date    string
	Content string
	Source  []byte
	URL     string
}

type ByDate []Post

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date > a[j].Date }

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

func writePost(post Post) {
	t, err := template.ParseFiles("templates/post.html")
	if err != nil {
		fmt.Printf("error %s", err)
	}
	fileName := fmt.Sprintf("public/%s.html", post.URL)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("error %s", err)
	}
	t.Execute(file, post)
}

func writePosts() []Post {
	posts := []Post{}
	files := getSources()
	for _, file := range files {
		post := parseSource(file)
		writePost(post)
		posts = append(posts, post)
	}
	return posts
}

func writeIndex(posts []Post) {
	sort.Sort(ByDate(posts))
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("error %s", err)
	}
	file, err := os.OpenFile("public/index.html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("error %s", err)
	}
	t.Execute(file, posts)
}

func main() {
	posts := writePosts()
	writeIndex(posts)
}
