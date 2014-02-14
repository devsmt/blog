package main

import (
	"bufio"
	"io"
	"net/http"
)

type DocumentParser interface {
	Parse(r io.Reader) (*Document, error)
}

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type Host string

func (h Host) Join(path string) string {
	sep := "/"
	if h[len(h)-1] != sep[0] {
		sep = ""
	}
	return string(h) + sep + path
}

type FileServer struct {
	host    Host
	dirfile string
	client  HttpClient
	parser  DocumentParser
}

func (fs *FileServer) Documents() ([]*Document, error) {
	rsp, err := fs.httpGet(fs.dirfile)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	docs := []*Document{}
	for s := bufio.NewScanner(rsp.Body); s.Scan(); {
		doc, err := fs.Get(s.Text())
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
func (fs *FileServer) httpGet(relpath string) (*http.Response, error) {
	githubPath := fs.host.Join(relpath)
	return fs.client.Get(githubPath)
}

func (fs *FileServer) Get(relpath string) (*Document, error) {
	rsp, err := fs.httpGet(relpath)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return fs.parser.Parse(rsp.Body)
}
