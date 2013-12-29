package main

import (
	"net/http"
	"github.com/weberc2/blog/github"
	"github.com/weberc2/blog/app"
	"log"
	"os"
)

func herokuport() string {
	return ":" + os.Getenv("PORT")
}

func main() {
	fetcher := github.NewFetcher("weberc2")
	a := app.New(fetcher)
	http.HandleFunc("/", a.HandleDocument)
	if err := http.ListenAndServe(herokuport(), nil); err != nil {
		log.Fatal(err)
	}
}
