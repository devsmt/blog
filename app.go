package main

import (
	"net/http"
	"text/template"
	"log"
)

type App struct {
	FileServer
	HomeTemplate, DocumentTemplate *template.Template
	Port string
}

// Home page handler: fetches documents from fileserver, reverses the order,
// and renders them into the HomeTemplate
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	docs, err := a.FileServer.Documents()
	if err != nil {
		httpErr(w, err, http.StatusInternalServerError)
		return
	}
	if err := a.HomeTemplate.Execute(w, reverse(docs)); err != nil {
		httpErr(w, err, http.StatusInternalServerError)
		return
	}
}

// Document page handler: fetches a document from fileserver (using the URL's
// path as the filepath for the fileserver) and rendering it into the
// DocumentTemplate
func (a *App) Document(w http.ResponseWriter, r *http.Request) {
	doc, err := a.FileServer.Get(r.URL.Path)
	if err != nil {
		httpErr(w, err, http.StatusInternalServerError)
		return
	}
	if err := a.DocumentTemplate.Execute(w, doc); err != nil {
		httpErr(w, err, http.StatusInternalServerError)
		return
	}
}

func httpErr(w http.ResponseWriter, err error, status int) {
	log.Println(err)
	http.Error(w, http.StatusText(status), status)
}

func reverse(docs []*Document) []*Document {
	reversed := make([]*Document, len(docs))
	for i:=0; i<len(docs); i++ {
		reversed[len(reversed)-1-i] = docs[i]
	}
	return reversed
}