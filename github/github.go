package github

import (
	"net/http"
	"io/ioutil"
	"fmt"
)

type Fetcher struct {
	root string
}

func (f *Fetcher) FetchDocument(path string) ([]byte, error) {
        path = "http://" + f.root + path
        resp, err := http.Get(path)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode > 299 || resp.StatusCode < 200 {
                return nil, fmt.Errorf("HTTP problem (path: %s): %s", path, resp.Status)
        }

        return ioutil.ReadAll(resp.Body)
}

func NewFetcher(username string) *Fetcher {
	return &Fetcher{root: fmt.Sprintf("%s.github.io", username)}
}
