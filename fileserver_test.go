package main

import (
	"log"
	"testing"
	"net/http"
)

type MockReadCloser struct {}
func (rc MockReadCloser) Read(_ []byte) (int, error) { return 0, nil }
func (rc MockReadCloser) Close() error { return nil }

type MockClient struct { Status int }
func (c *MockClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, MockReadCloser{})
	if err != nil {
		log.Fatal("Unexpected err:", err)
	}
	return &http.Response {
		Status: http.StatusText(c.Status),
		StatusCode: c.Status,
		Body: MockReadCloser{},
		Request: req,
	}, nil
}

type Fixture struct {
	FileServer
	Client *MockClient
}

func fixture() *Fixture {
	client := &MockClient{}
	return &Fixture {
		Client: client,
		FileServer: FileServer {
			client: client,
			parser: &MetadataParser{},
		},
	}
}

// When the HTTP Client returns 404 during a Get(), expect the resultant error
// satisfies IsNotExist()
func TestGet_HttpClientReturns404(t *testing.T) {
	fs := fixture()
	fs.Client.Status = http.StatusNotFound
	_, err := fs.Get("")
	if !fs.IsNotExist(err) {
		t.Fatalf("Expected IsNotExist(err) to be true; err: %v", err)
	}
}
