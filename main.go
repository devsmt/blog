package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

func parseTemplate(filepath ...string) *template.Template {
	return template.Must(template.ParseFiles(filepath...))
}

func main() {
	log.SetFlags(log.Lshortfile)
	app := App{
		DocumentStore: &FileServer{
			//host: "http://localhost:8000",
			host:    "http://weberc2.github.io/",
			dirfile: "dirfile",
			client:  http.DefaultClient,
			parser:  new(MetadataParser),
		},
		HomeTemplate:     parseTemplate("home.html", "google-analytics.html", "css.html", "base.html"),
		DocumentTemplate: parseTemplate("document.html", "google-analytics.html", "css.html", "base.html", "disqus.html"),
		Port:             ":" + os.Getenv("PORT"),
		PageSize:         10,
	}

	if err := app.Run(); err != nil {
		log.Println("Error running app:", err)
		return
	}
}
