package main

import (
	"log"
	"net/http"
	"text/template"
)

func parseTemplate(filepath ...string) *template.Template {
	return template.Must(template.ParseFiles(filepath...))
}

func main() {
	app := App{
		DocumentStore: &FileServer{
			//host: "http://localhost:8000",
			host:    "http://weberc2.github.io/",
			dirfile: "dirfile",
			client:  http.DefaultClient,
			parser:  new(MetadataParser),
		},
		HomeTemplate:     parseTemplate("home.html", "base.html"),
		DocumentTemplate: parseTemplate("document.html", "base.html"),
		Port:             ":8080",
		PageSize:         10,
	}

	if err := app.Run(); err != nil {
		log.Println(err)
		return
	}
}
