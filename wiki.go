package main

import (
	"fmt"
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
    p, _ := loadPage(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}


func editHandler(w http.ResponseWriter, r *http.Request) {
	
	// extract  page title from r.URL.Path
    title := r.URL.Path[len("/edit/"):]

    // load page data and format with html 
    p, err := loadPage(title)

    // handle error
    if err != nil {
        p = &Page{Title: title}
    }
    
    fmt.Fprintf(w, "<h1>Editing %s</h1>"+
        "<form action=\"/save/%s\" method=\"POST\">"+
        "<textarea name=\"body\">%s</textarea><br>"+
        "<input type=\"submit\" value=\"Save\">"+
        "</form>",
        p.Title, p.Title, p.Body)
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

