package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

// Post Data
type Post struct {
	Content string
	Meta    *Meta
}

// Meta Data
type Meta struct {
	Title string
	Tags  []string
	Date  string
	Slug  string `yaml:"permalink"`
	Draft bool
}

// ByDate use for post sort
type ByDate []Post

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Meta.Date > a[j].Meta.Date }

func getSources() []string {
	files, _ := filepath.Glob("srcs/*.md")
	return files
}

func renderMarkdown(source []byte) string {
	content := string(blackfriday.Run(source))
	return content
}

func parseSource(fileName string) Post {
	post := Post{}
	sources, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(sources), "\n")
	metaLoc := [2]int{}
	SEP := "---"
	if lines[0] == SEP {
		metaLoc[0] = 1
	}
	for k, v := range lines[1:] {
		if v == SEP {
			if metaLoc[1] != 0 {
				break
			}
			metaLoc[1] = k + 1
		}
	}
	meta := lines[metaLoc[0]:metaLoc[1]]
	metaSource := []byte(strings.Join(meta, "\n"))
	source := []byte(strings.Join(lines[metaLoc[1]:len(lines)], "\n"))

	err := yaml.Unmarshal(metaSource, &post.Meta)
	if err != nil {
		fmt.Printf("error %s", err)
	}
	post.Content = renderMarkdown(source)
	return post
}

func writePost(post Post) {
	t, err := template.ParseFiles("templates/post.html")
	if err != nil {
		fmt.Printf("error %s", err)
	}
	fileName := fmt.Sprintf("public/%s.html", post.Meta.Slug)
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

func handler(w http.ResponseWriter, r *http.Request) {
	var validPath = regexp.MustCompile("^/([a-z0-9-]+)$")
	postURL := validPath.FindStringSubmatch(r.URL.Path)
	filePath := "public/index.html"

	if postURL != nil {
		filePath = fmt.Sprintf("public/%s.html", postURL[1])
	}
	log.Println(filePath)
	t, _ := template.ParseFiles(filePath)
	err := t.Execute(w, nil)
	if err != nil {
		fmt.Printf("error %s", err)
	}
}

func main() {
	posts := writePosts()
	writeIndex(posts)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
