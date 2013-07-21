package main

import (
	"html/template"
	"net/http"
)

var (
	viewTemplate = template.Must(template.ParseFiles("tmpl/base.html", "tmpl/view.html"))
	editTemplate = template.Must(template.ParseFiles("tmpl/base.html", "tmpl/edit.html"))
)

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, p *Page) {
	err := tmpl.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
