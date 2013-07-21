package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

const (
	viewPath, viewTemplate = "/view/", "view.html"
	editPath, editTemplate = "/edit/", "edit.html"
	savePath               = "/save/"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{title, body}, err
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(viewPath):]
	page, _ := loadPage(title)

	t, _ := template.ParseFiles(viewTemplate)
	t.Execute(w, page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(editPath):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{title, nil}
	}

	t, _ := template.ParseFiles(editTemplate)
	t.Execute(w, page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	page := &Page{"TestPage", []byte("This is a sample page.")}
	page.save()

	http.HandleFunc(viewPath, viewHandler)
	http.HandleFunc(editPath, editHandler)
	http.HandleFunc(savePath, saveHandler)
	http.ListenAndServe("localhost:8080", nil)
}
