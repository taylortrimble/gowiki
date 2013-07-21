package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

const (
	viewPath = "/view/"
	editPath = "/edit/"
	savePath = "/save/"
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(viewPath):]
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, editPath+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(editPath):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{title, nil}
	}

	renderTemplate(w, "edit.html", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len(savePath):]
	body := r.FormValue("body")

	p := &Page{title, []byte(body)}
	p.save()

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	page := &Page{"TestPage", []byte("This is a sample page.")}
	page.save()

	http.HandleFunc(viewPath, viewHandler)
	http.HandleFunc(editPath, editHandler)
	http.HandleFunc(savePath, saveHandler)
	http.ListenAndServe("localhost:8080", nil)
}
