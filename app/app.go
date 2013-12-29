package app
import (
	"net/http"
	"log"
	"fmt"
	"github.com/russross/blackfriday"
)

type DocumentFetcher interface {
	FetchDocument(path string) ([]byte, error)
}

type App struct {
	fetcher DocumentFetcher
}

func New(fetcher DocumentFetcher) *App {
	return &App{fetcher: fetcher}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "404 Not Found")
}

func (a *App) HandleDocument(w http.ResponseWriter, r *http.Request) {
	data, err := a.fetcher.FetchDocument(r.URL.Path)
	if err != nil {
		NotFound(w,r)
		log.Println(err)
	}

	data = blackfriday.MarkdownCommon(data)

	if _, err := w.Write(data); err != nil {
		log.Println(err)
	}
}
