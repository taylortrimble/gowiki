package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	viewPath = "/view/"
	editPath = "/edit/"
	savePath = "/save/"
)

var (
	templates      = template.Must(template.ParseFiles("view.html", "edit.html"))
	titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
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
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, path string, r *http.Request) (title string, err error) {
	title = r.URL.Path[len(path):]
	if !titleValidator.MatchString(title) {
		http.NotFound(w, r)
		err = errors.New("Invalid page title.")
	}
	return
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, viewPath, r)
	if err != nil {
		return
	}

	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, editPath+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, editPath, r)
	if err != nil {
		return
	}

	page, err := loadPage(title)
	if err != nil {
		page = &Page{title, nil}
	}

	renderTemplate(w, "edit.html", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, savePath, r)
	if err != nil {
		return
	}

	body := r.FormValue("body")

	p := &Page{title, []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
