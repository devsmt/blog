package main

import (
	"fmt"
	"net/http"
	"log"
	"io"
)

type Template interface {
	Execute(w io.Writer, data interface{}) error
}

type DocumentStore interface {
	Get(addr string) (*Document, error)
	Documents() ([]*Document, error)
	IsNotExist(err error) bool
}

type App struct {
	DocumentStore
	HomeTemplate, DocumentTemplate Template
	Port string
}

// Home page handler: fetches documents from fileserver, reverses the order,
// and renders them into the HomeTemplate
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	docs, err := a.DocumentStore.Documents()
	if a.DocumentStore.IsNotExist(err) {
		httpErr(w, err, http.StatusNotFound)
	} else if err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	} else if err := a.HomeTemplate.Execute(w, reverse(docs)); err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	}
}

// Document page handler: fetches a document from fileserver (using the URL's
// path as the filepath for the fileserver) and rendering it into the
// DocumentTemplate
func (a *App) Document(w http.ResponseWriter, r *http.Request) {
	doc, err := a.DocumentStore.Get(r.URL.Path)
	if a.DocumentStore.IsNotExist(err) {
		httpErr(w, err, http.StatusNotFound);
	} else if err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	} else if err := a.DocumentTemplate.Execute(w, doc); err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	}
}

func (a *App) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			a.Home(w, r)
		default:
			a.Document(w, r)
		}
	})

	return http.ListenAndServe(a.Port, nil)
}

/* Helpers */

func httpErrText(status int) string {
	return fmt.Sprintf("%d ", status) + http.StatusText(status)
}

func httpErr(w http.ResponseWriter, err error, status int) {
	log.Println(err)
	http.Error(w, httpErrText(status), status)
}

func reverse(docs []*Document) []*Document {
	reversed := make([]*Document, len(docs))
	for i:=0; i<len(docs); i++ {
		reversed[len(reversed)-1-i] = docs[i]
	}
	return reversed
}
