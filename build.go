package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

// Post Data
type Post struct {
	Content string
	Excerpt string
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

// Tag Page Data
type Tag struct {
	Name  string
	Posts []Post
}

// Table is template Data
type Table map[string]interface{}

// ByDate use for post sort
type ByDate []Post

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Meta.Date > a[j].Meta.Date }

const (
	// SEP split post Meta and Content
	SEP = "---"
	// CURDIR is Current directory
	CURDIR = "."
	// UPDIR is Upper level directory
	UPDIR = ".."
	// TPLDIR is template directory
	TPLDIR = "./templates/"
	// PUBDIR is public directory
	PUBDIR = "public"
)

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
	content := strings.Join(lines[metaLoc[1]+1:len(lines)], "\n")
	source := []byte(content)

	err := yaml.Unmarshal(metaSource, &post.Meta)
	if err != nil {
		fmt.Printf("error %s", err)
	}
	if len(post.Excerpt) > 100 {
		post.Excerpt = content[:100]
	} else {
		post.Excerpt = content
	}
	post.Content = renderMarkdown(source)
	return post
}

func writePost(post Post) {
	fileName := fmt.Sprintf(PUBDIR+"/%s.html", post.Meta.Slug)
	err := renderTemplate(fileName, "post.html", Table{"Post": post, "Prefix": CURDIR})
	if err != nil {
		fmt.Printf("error %s", err)
	}
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

func writeTagPage(tags map[string][]Post) {
	for tag, post := range tags {
		fileName := fmt.Sprintf(PUBDIR+"/tag/%s.html", tag)
		err := renderTemplate(fileName, "tag.html", Table{"Tag": Tag{tag, post}, "Prefix": UPDIR})
		if err != nil {
			fmt.Printf("error %s", err)
		}
	}
}

func writeTagsIndex(tags map[string][]Post) {
	err := renderTemplate(PUBDIR+"/tags.html", "tags.html", Table{"Tags": tags, "Prefix": CURDIR})
	if err != nil {
		fmt.Printf("error %s", err)
	}
}

func writeTags(posts []Post) {

	tags := make(map[string][]Post, 0)
	for _, post := range posts {
		for _, tag := range post.Meta.Tags {
			tags[tag] = append(tags[tag], post)
		}
	}
	writeTagPage(tags)
	writeTagsIndex(tags)
}

func writeIndex(posts []Post) {
	sort.Sort(ByDate(posts))
	err := renderTemplate(PUBDIR+"/index.html", "index.html", Table{"Post": posts, "Prefix": CURDIR})
	if err != nil {
		fmt.Printf("error %s", err)
	}
}

var templates map[string]*template.Template

func parseTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layouts, err := filepath.Glob(TPLDIR + "layouts/*.html")
	if err != nil {
		log.Fatal(err)
	}
	partials, err := filepath.Glob(TPLDIR + "partials/*.html")
	if err != nil {
		log.Fatal(err)
	}
	for _, layout := range layouts {
		files := append(partials, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
}

func renderTemplate(filePath string, tempName string, data interface{}) error {
	tmpl, ok := templates[tempName]
	if !ok {
		return fmt.Errorf("The template %s does not exist", tempName)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("%s create fail", filePath)
	}
	return tmpl.ExecuteTemplate(file, tempName, data)
}
func cleanDir(dst string) error {
	cmd := exec.Command("rm", "-rf", dst)
	err := cmd.Run()
	return err
}

func createDir(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func copyDir(src string, dst string) error {
	info, err := os.Stat(src)
	if err != nil || !info.IsDir() {
		return err
	}
	cmd := exec.Command("cp", "-rf", src, dst)
	err = cmd.Run()
	return err
}

func createPaths() {
	err := cleanDir(PUBDIR)
	if err != nil {
		log.Fatal(err)
	}
	err = createDir(PUBDIR)
	if err != nil {
		log.Fatal(err)
	}
	err = copyDir(TPLDIR+"static", PUBDIR)
	if err != nil {
		log.Fatal(err)
	}
	err = createDir(PUBDIR + "/tag")
	if err != nil {
		log.Fatal(err)
	}
}

func generate() {
	parseTemplates()
	createPaths()
	posts := writePosts()
	writeIndex(posts)
	writeTags(posts)
}

func server() {
	http.Handle("/", http.FileServer(http.Dir(PUBDIR)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var (
	serve bool
	gene  bool
)

func init() {
	flag.BoolVar(&serve, "s", false, "server on 8080")
	flag.BoolVar(&gene, "g", false, "clean and generate")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, `go-static-blog version: 1.0.0
Usage: go-static-blog [-g generate] [-s server] 

Options:
`)
	flag.PrintDefaults()
}

func main() {
	if gene {
		generate()
	}
	if serve {
		server()
	}
}
