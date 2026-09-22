package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := "templates/" + title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

var validPath = regexp.MustCompile("^/(|home|homelab|projects)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			//http.NotFound(w, r)
			fn(w, r, "home")
			return
		}
		title := m[1]
		if title == "" {
			title = "home"
		}
		fn(w, r, title)
	}
}

var templates = template.Must(template.ParseFiles("templates/home.html", "templates/homelab.html", "templates/projects.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch title {
	case "homelab":
		renderTemplate(w, "homelab", p)
	case "home":
		renderTemplate(w, "home", p)
	case "projects":
		renderTemplate(w, "projects", p)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/", makeHandler(viewHandler))
	http.HandleFunc("/home", makeHandler(viewHandler))
	http.HandleFunc("/homelab", makeHandler(viewHandler))
	http.HandleFunc("/projects", makeHandler(viewHandler))

	// CSS shit
	dir := http.Dir("./static")
	handler := http.StripPrefix("/static/", http.FileServer(dir))
	http.Handle("/static/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
