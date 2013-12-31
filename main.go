package main

import (
	"io"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"bufio"
)

func mustEnvVar(varname string) string {
	env := os.Getenv(varname)
	if env == "" {
		log.Fatal("%s not set", varname)
	}
	return env
}

func url(relpath string) string {
	return fmt.Sprintf("http://%s.github.io/%s", GITHUB_USER, relpath)
}

var (
	GITHUB_USER = mustEnvVar("GITHUB_USER")
	HEROKU_PORT = ":" + mustEnvVar("PORT")
	DIRECTORY_FILE = url("dirfile")
	TITLE = "weberc2"
)

func init() {
	header := fmt.Sprintf(`<html><body><h1 class="main-header"><a href="/">%s</a></h1>`, TITLE)
	SetHeader([]byte(header))
	SetFooter([]byte(`</body></html>`))
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

		Write(w, data)
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

		WriteHome(w, docs)
	})

	if err := http.ListenAndServe(HEROKU_PORT, r); err != nil {
		log.Fatal(err)
	}
}
