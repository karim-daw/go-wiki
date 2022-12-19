package main

import (
	"errors"
	"fmt"
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

// write to a file given a page struct, returns error
func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}

// global variable that holds all templates
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// global variable to store validation expression
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


// function that uses the validPath expression to validate path and extract the page title
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    stringMatch := validPath.FindStringSubmatch(r.URL.Path)
    if stringMatch == nil {
        http.NotFound(w, r)
        return "", errors.New("invalid Page Title")
    }
    return stringMatch[2], nil // The title is the second subexpression.
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

	// extract  page title with getTitle function
    title, err := getTitle(w, r)
    if err != nil {
        return
    }

	// load page data and format with html 
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}


func editHandler(w http.ResponseWriter, r *http.Request) {
	
	// extract  page title with getTitle function
    title, err := getTitle(w, r)
    if err != nil {
        return
    }

    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}


func saveHandler(w http.ResponseWriter, r *http.Request){

	// extract  page title with getTitle function
    title, err := getTitle(w, r)
    if err != nil {
        return
    }

    body := r.FormValue("body")

    // convert string from body to bytes
    p := &Page{Title: title, Body: []byte(body)}

    err = p.save()

    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
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

