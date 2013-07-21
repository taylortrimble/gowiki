package main

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
)

const (
	viewPath = "/view/"
	editPath = "/edit/"
	savePath = "/save/"
)

var (
	templates      = template.Must(template.ParseFiles("tmpl/view.html", "tmpl/edit.html"))
	titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
)

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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, editPath+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	page, err := loadPage(title)
	if err != nil {
		page = &Page{title, nil}
	}

	renderTemplate(w, "edit.html", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")

	p := &Page{title, []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(w, savePath, r)
		if err != nil {
			return
		}
		fn(w, r, title)
	}
}

func main() {
	page := &Page{"TestPage", []byte("This is a sample page.")}
	page.save()

	http.HandleFunc(viewPath, makeHandler(viewHandler))
	http.HandleFunc(editPath, makeHandler(editHandler))
	http.HandleFunc(savePath, makeHandler(saveHandler))
	http.ListenAndServe("localhost:8080", nil)
}
