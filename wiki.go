package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const viewPath = "/view/"

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
	fmt.Fprintf(w, "<h1>%s</h1><body>%s</body>", page.Title, page.Body)
}

func main() {
	page := &Page{"TestPage", []byte("This is a sample page.")}
	page.save()

	http.HandleFunc(viewPath, viewHandler)
	http.ListenAndServe("localhost:8080", nil)
}
