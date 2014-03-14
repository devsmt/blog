package testsupport

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"io"
)

type MockReadCloser struct {
	io.Reader
}

func NewMockReadCloser(data string) *MockReadCloser {
	return &MockReadCloser{Reader: strings.NewReader(data)}
}

func (rc *MockReadCloser) Close() error { return nil }

type MockClient struct {
	Status int
	files  map[string]string
}

func NewMockClient() *MockClient {
	return &MockClient {
		Status: http.StatusOK,
		files:  map[string]string{},
	}
}

func (rc *MockClient) AddFile(name, contents string) {
	rc.files[name] = contents
}

func (rc *MockClient) AppendToFile(name, text string) {
	rc.files[name] += text
}

func (c *MockClient) Get(url string) (*http.Response, error) {
	println("Request for", url)
	req, err := http.NewRequest("GET", url, NewMockReadCloser(""))
	if err != nil {
		log.Fatal("Unexpected err:", err)
	}

	return &http.Response{
		Status:     fmt.Sprintf("%d ", c.Status) + http.StatusText(c.Status),
		StatusCode: c.Status,
		Body:       NewMockReadCloser(""),
		Request:    req,
	}, nil
}
