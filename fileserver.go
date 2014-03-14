package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

type FileServer struct {
	host    host           // URL of host
	dirfile string         // relative path of directory file
	client  HttpClient     // an http client for communicating w/ server
	parser  DocumentParser // a document parser for interpreting the remote file
}

// Reads the list of published documents from the server's dirfile,
// gets those messages, and returns them. Will return (nil, err) as
// soon as an error is encountered
func (fs *FileServer) Documents(start, end int) ([]*Document, error) {
	rsp, err := fs.httpGet(fs.dirfile)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	docs := []*Document{}
	println("Got here")
	for i, s := 0, bufio.NewScanner(rsp.Body); s.Scan() || i >= end; i++ {
		println("And here")
		fmt.Printf("Line %d of dirfile: %v\n", i, s.Text())
		// skip to the beginning
		if i < start {
			continue
		}

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

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(rsp.Status)
	}

	doc, err := fs.parser.Parse(rsp.Body)
	if err != nil {
		return nil, err
	}
	doc.Metadata["Path"] = relpath
	return doc, nil
}

// Return whether `err` was caused by attempting to access an unknown or invalid
// path/address
func (fs *FileServer) IsNotExist(err error) bool {
	statusText := http.StatusText(http.StatusNotFound)
	notFoundErrMsg := fmt.Sprintf("%d ", http.StatusNotFound) + statusText
	return err != nil && err.Error() == notFoundErrMsg
}

/* Helpers */

type DocumentParser interface {
	Parse(r io.Reader) (*Document, error)
}

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type host string

func (h host) Join(path string) string {
	if len(h) < 1 {
		return path
	}

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
