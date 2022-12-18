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
func (p *Page) save() error {
    filename := p.Title + ".txt"
    return os.WriteFile(filename, p.Body, 0600)
}