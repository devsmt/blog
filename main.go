package main

import (
	"bufio"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	GITHUB_USER    = "weberc2"
	HEROKU_PORT    = ":" + mustEnvVar("PORT")
	DIRECTORY_FILE = url("dirfile")
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	rsp, err := http.Get(DIRECTORY_FILE)
	if err != nil {
		internalServerErr(w, err)
		return
	}
	defer rsp.Body.Close()

	paths := []string{}
	for s := bufio.NewScanner(rsp.Body); s.Scan(); {
		paths = append(paths, s.Text())
	}

	docs := []*document{}
	for i := len(paths) - 1; i >= 0; i-- {
		path := paths[i]
		rsp, err := http.Get(url(path))
		if err != nil {
			log.Println("Error getting post content:", err)
			continue
		}
		defer rsp.Body.Close()

		doc, err := parseDoc(path, rsp.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		docs = append(docs, doc)
	}

	HOME_TEMPLATE.Execute(w, docs)
}

func documentHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

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

	DOC_TEMPLATE.Execute(w, string(blackfriday.MarkdownCommon(data)))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			homeHandler(w, r)
		default:
			documentHandler(w, r)
		}
	})

	if err := http.ListenAndServe(HEROKU_PORT, nil); err != nil {
		log.Fatal(err)
	}
}
