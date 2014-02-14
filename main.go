package main

import (
	"net/http"
	"fmt"
)

type App struct {
	FileServer
	Templates struct {
		Home, Document interface {}
	}
	Port string
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	docs, err := a.FileServer.Documents()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(w, docs)
}

func (a *App) Document(w http.ResponseWriter, r *http.Request) {
	doc, err := a.FileServer.Get(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, doc)
}

func main() {
	app := App {
		FileServer: FileServer{
			host: "http://weberc2.github.io/",
			dirfile: "dirfile",
			client: http.DefaultClient,
			parser: new(FakeDocParser),
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			app.Home(w, r)
		default:
			app.Document(w, r)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
}
