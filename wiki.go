package main

import (
	"os"
)

// A wiki has a title and a body (the page content).
// The Body element is a []byte rather than string because that is expected by the io libraries
type Page struct {
    Title string
    Body  []byte
}

// This is a method named save that takes as its receiver
// p, a pointer to Page . It takes no parameters, and returns a value of type error
func (page *Page) save() error {
    filename := page.Title + ".txt"
    return os.WriteFile(filename, page.Body, 0600)
}

// method loads a page if it exists, if not, returns error
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func main() {
    page1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    page1.save()
    // page2, _ := loadPage("TestPage")
    // fmt.Println(string(page2.Body))
}