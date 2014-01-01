package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"bufio"
	"text/template"
)

func mustEnvVar(varname string) string {
	env := os.Getenv(varname)
	if env == "" {
		log.Fatalf("%s not set", varname)
	}
	return env
}

func url(relpath string) string {
	return fmt.Sprintf("http://%s.github.io/%s", GITHUB_USER, relpath)
}

var (
	GITHUB_USER = "weberc2"
	HEROKU_PORT = ":" + mustEnvVar("PORT")
	DIRECTORY_FILE = url("dirfile")
	TITLE = "weberc2"
	HOME_RENDERER = NewMarkdownRenderer(fetchTemplate(url("tmpl/home.html")))
	PAGE_RENDERER = NewMarkdownRenderer(fetchTemplate(url("tmpl/page.html")))
)

func fetch(filename string) ([]byte, error) {
	rsp, err := http.Get(filename)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 299 || rsp.StatusCode < 200 {
		return nil, fmt.Errorf("HTTP problem fetching %s: %s", filename, rsp.Status)
	}

	return ioutil.ReadAll(rsp.Body)
}

func fetchTemplate(filename string) *template.Template {
	data, err := fetch(filename)
	if err != nil {
		log.Fatal(err)
	}

	return template.Must(template.New(filename).Parse(string(data)))
}

func main() {
	r := mux.NewRouter()

	// document handler
	r.HandleFunc("/{path}", func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)["path"]

		resp, err := http.Get(url(path))
		if err != nil {
			internalServerErr(w, err)
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			internalServerErr(w, err)
			return
		}

		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			w.WriteHeader(resp.StatusCode)
			fmt.Fprintln(w, "HTTP error:", resp.Status)
			return
		}

		if err := PAGE_RENDERER.Write(w, data); err != nil {
			log.Println(err)
		}
	})

	// home handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rsp, err := http.Get(DIRECTORY_FILE)
		if err != nil {
			internalServerErr(w, err)
			return
		}
		defer rsp.Body.Close()

		docs := []*document{}
		for s := bufio.NewScanner(rsp.Body); s.Scan(); {
			path := s.Text()
			rsp, err := http.Get(url(path))
			if err != nil {
				log.Println("Error getting post content:", err)
				continue
			}
			defer rsp.Body.Close()

			docs = append(docs, parseDoc(path, rsp.Body))
		}

		if err := HOME_RENDERER.WriteHome(w, docs); err != nil {
			log.Println(err)
		}
	})

	if err := http.ListenAndServe(HEROKU_PORT, r); err != nil {
		log.Fatal(err)
	}
}
