package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(reading|posts|)(/*)?")

func main() {
	http.HandleFunc("/", makeHandler(viewHandler))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type Page struct {
	View string
	Title string
	Body []byte // switch this to parses markdown
	Posts []Post
}

func setPage(title string) (*Page, error) {
	view := getView(title)
	filename := "./pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename) // switch this to markdown parse
	if err != nil {
		return nil, err
	}
	var posts []Post

	if title == "posts" {
		posts = setPosts()
	}

	return &Page{View: view, Title: title, Body: body, Posts: posts}, nil
}

type Post struct {
	Title string
	URL string
}

func setPosts() []Post {
	var posts []Post

	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		item := Post{Title: "Foo", URL: "./posts/" + file.Name()}
		posts = append(posts, item)
	}
	return posts
}

func makeHandler(fn func (w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var title string
		if r.URL.Path == "/" {
			title = "home"
		} else {
			m := validPath.FindStringSubmatch(r.URL.Path)
			title = m[1] // The title is the first subexpression.
		}
		fn(w, r, title)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := setPage(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, p)
}

func renderTemplate(w http.ResponseWriter, p *Page) {
	t, err := template.ParseFiles(
		"./views/" + p.View + ".html",
		"./views/head.html",
		"./views/scripts.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getView(title string) string {
	switch title {
		case "home":
			return "home"
		case "posts":
			return "page"
		default:
			return "post"
	}
}