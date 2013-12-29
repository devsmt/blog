package main

import (
	"net/http"
	"github.com/weberc2/blog/github"
	"github.com/weberc2/blog/app"
	"log"
)

func main() {
	fetcher := github.NewFetcher("weberc2")
	a := app.New(fetcher)
	http.HandleFunc("/", a.HandleDocument)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
