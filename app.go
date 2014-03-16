package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"strconv"
)

const pagePrefix = "/page/"

type Template interface {
	ExecuteTemplate(w io.Writer, name string, data interface{}) error
}

type DocumentStore interface {
	Get(addr string) (*Document, error)
	Documents(start, end int) ([]*Document, error)
	IsNotExist(err error) bool
}

type App struct {
	DocumentStore
	HomeTemplate, DocumentTemplate Template
	Port                           string
	PageSize                       int
}

func (a *App) Page(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, pagePrefix) {
		log.Fatalf("Expected path begins with '%s', got '%s'", pagePrefix, r.URL.Path)
	}
	pageString := r.URL.Path[len(pagePrefix):]
	pageNumber, err := strconv.Atoi(pageString)
	if err != nil {
		httpErr(w, err, http.StatusNotFound)
		return
	}
	a.pageHandler(w, pageNumber)
}

func (a *App) pageHandler(w http.ResponseWriter, pageNumber int) {
	start := pageNumber * a.PageSize
	end := start + a.PageSize
	docs, err := a.DocumentStore.Documents(start, end)
	if a.DocumentStore.IsNotExist(err) {
		httpErr(w, err, http.StatusNotFound)
	} else if err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	}

	homePage := struct {
		Docs []*Document
		PreviousPage, NextPage string
	} {
		Docs: docs,
	}

	previousPage := pageNumber - 1
	if previousPage >= 0 {
		homePage.PreviousPage = strconv.Itoa(previousPage)
	}

	nextPage := pageNumber + 1
//	 if nextPage <= a.DocumentStore.DocumentCount() / a.PageSize {
	homePage.NextPage = strconv.Itoa(nextPage)
//	 }

	if err := a.HomeTemplate.ExecuteTemplate(w, "base", homePage); err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	}
}

// Home page handler: fetches documents from fileserver, reverses the order,
// and renders them into the HomeTemplate
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	a.pageHandler(w, 0)
}

// Document page handler: fetches a document from fileserver (using the URL's
// path as the filepath for the fileserver) and rendering it into the
// DocumentTemplate
func (a *App) Document(w http.ResponseWriter, r *http.Request) {
	doc, err := a.DocumentStore.Get(r.URL.Path)
	if a.DocumentStore.IsNotExist(err) {
		httpErr(w, err, http.StatusNotFound)
	} else if err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	} else if err := a.DocumentTemplate.ExecuteTemplate(w, "base", doc); err != nil {
		httpErr(w, err, http.StatusInternalServerError)
	}
}

func (a *App) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			a.Home(w, r)
		} else if strings.HasPrefix(path, pagePrefix) {
			a.Page(w, r)
		} else {
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
