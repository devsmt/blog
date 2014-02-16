package main

import (
	"net/http"
	"fmt"
	"text/template"
)

func parseTemplate(filepath string) *template.Template {
	return template.Must(template.ParseFiles(filepath))
}

func main() {
	app := App {
		FileServer: FileServer{
			//host: "http://localhost:3000",
			host: "http://weberc2.github.io/",
			dirfile: "dirfile",
			client: http.DefaultClient,
			parser: new(MetadataParser),
		},
		HomeTemplate: parseTemplate("home.html"),
		DocumentTemplate: parseTemplate("document.html"),
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
