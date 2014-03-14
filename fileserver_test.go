package main

import (
	"io"
	"log"
	"net/http"
	"testing"
	"strconv"
)

type FakeReadCloser struct {
	data string
}

func (rc *FakeReadCloser) Read(p []byte) (int, error) {
	p = []byte(rc.data)
	return len(rc.data), nil
}
func (rc *FakeReadCloser) Close() error { return nil }

type MockReadCloser struct{}

func (rc MockReadCloser) Read(_ []byte) (int, error) { return 0, nil }
func (rc *MockReadCloser) Close() error               { return nil }

type MockClient struct {
	Status int
	Pages  map[string]string
}

func (c *MockClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, MockReadCloser{})
	if err != nil {
		log.Fatal("Unexpected err:", err)
	}

	var rc io.ReadCloser
	if page, found := c.Pages[url]; found {
		rc = FakeReadCloser{data:page}
	} else {
		rc = MockReadCloser{}
	}

	return &http.Response{
		Status:     httpErrText(c.Status),
		StatusCode: c.Status,
		Body:       rc,
		Request:    req,
	}, nil
}

type Fixture struct {
	FileServer
	Client *MockClient
}

func fixture() *Fixture {
	client := &MockClient{Status: http.StatusOK}
	return &Fixture{
		Client: client,
		FileServer: FileServer{
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

func TestDocuments(t *testing.T) {
	fs := fixture()
	// Given 5 documents in fileserver
	pages := map[string]string{
		"1" : "first",
		"2" : "second",
		"3" : "third",
		"4" : "fourth",
		"5" : "fifth",
	}
	fs.Client.Pages	= pages
	// When documents 0-5 are requested
	docs, err := fs.Documents(0, 5)
	if err != nil {
		t.Fatal("Unexpected err: ", err)
	}

	// Expect all 5 documents are returned
	if n := len(docs); n != 5 {
		t.Fatalf("Expected 5 documents, found %d", n)
	}

	for i, doc := range docs {
		url := strconv.Itoa(i+1)
		exp := pages[url]
		act := string(doc.Text)
		if exp != act {
			t.Fatalf("Expected %s; got %s", exp, act)
		}
	}
}

func TestDocuments_WhenEndExceedsNumberOfDocsInDirectoryFile(t *testing.T) {
	// Given 5 documents in directory file
	// When documents 0-6 are requested
	// Expect that 5 documents are returned
	t.Fatal("Implement me")
}

func TestDocuments_WhenEmpty(t *testing.T) {
	// Given no documents in directory file
	// When documents 0-1 are requested
	// Expect that no documents are returned
	t.Fatal("Implement me")
}
