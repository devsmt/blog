package main

import (
	"bufio"
	"io"
	"net/http"
)

type FileServer struct {
	host    host           // URL of host
	dirfile string         // relative path of directory file
	client  httpClient     // an http client for communicating w/ server
	parser  documentParser // a document parser for interpreting the remote file
}

// Reads the list of published documents from the server's dirfile,
// gets those messages, and returns them. Will return (nil, err) as
// soon as an error is encountered
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

// Gets a document from the fileserver using `relpath` parameter as
// the file's address; TODO: What happens if the server returns
// a non-200 value (e.g., 404)?
func (fs *FileServer) Get(relpath string) (*Document, error) {
	rsp, err := fs.httpGet(relpath)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	doc, err := fs.parser.Parse(rsp.Body)
	if err != nil {
		return nil, err
	}
	doc.Metadata["Path"] = relpath
	return doc, nil
}

/* Helpers */

type documentParser interface {
	Parse(r io.Reader) (*Document, error)
}

type httpClient interface {
	Get(url string) (*http.Response, error)
}

type host string

func (h host) Join(path string) string {
	sep := ""
	if h[len(h)-1] != '/' {
		sep = "/"
	}
	return string(h) + sep + path
}

func (fs *FileServer) httpGet(relpath string) (*http.Response, error) {
	githubPath := fs.host.Join(relpath)
	return fs.client.Get(githubPath)
}
