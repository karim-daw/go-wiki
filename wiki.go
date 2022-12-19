package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)


type Page struct {
    Title string
    Body  []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}

// loadPage will return a pointer to a page struct based on
// the pages title string
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}


// handle root
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}


// handle view end point
func viewHandler(w http.ResponseWriter, r *http.Request) {

	// extract  page title from r.URL.Path
    title := r.URL.Path[len("/view/"):]

	// load page data and format with html 
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}


func editHandler(w http.ResponseWriter, r *http.Request) {
	
	// extract  page title from r.URL.Path
    // load page data and format with html 
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}


func saveHandler(w http.ResponseWriter, r *http.Request){

    // get string body from form value and turn into page struct
    title := r.URL.Path[len("/save/"):]
    body := r.FormValue("body")

    // convert string from body to bytes
    p := &Page{Title: title, Body: []byte(body)}
    p.save()
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
}


func main() {

	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    p1.save()

    http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
    http.HandleFunc("/save/", saveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

