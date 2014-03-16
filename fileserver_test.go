package main

import (
	"net/http"
	"testing"
	. "github.com/weberc2/blog/testsupport"
	"strconv"
)

type Fixture struct {
	FileServer
	Client *MockClient
}

func (f *Fixture) NewFile(name, content string) {
	f.Client.AppendLineToFile(f.dirfile, name)
	f.Client.AddFile(name, content)
}

func fixture() *Fixture {
	client := NewMockClient()
	return &Fixture{
		Client: client,
		FileServer: FileServer{
			client: client,
			parser: &MetadataParser{},
			dirfile: "dirfile",
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
	// Given 5 documents in fileserver
	f := fixture()
	files := map[string]string {
		"1" : "first",
		"2" : "second",
		"3" : "third",
		"4" : "fourth",
		"5" : "fifth",
	}

	for name, content := range files {
		f.NewFile(name, content)
	}

	// When documents 0-5 are requested
	docs, err := f.Documents(0, 5)
	if err != nil {
		t.Fatalf("Unexpected err: ", err)
	}

	// Expect all 5 documents are returned
	if n := len(docs); n != 5 {
		t.Fatalf("Expected 5; Got %d", n)
	}

	for i, doc := range docs {
		exp := files[strconv.Itoa(i+1)]
		act := string(doc.Text)
		if err := Compare(exp, act); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDocuments_WhenEndExceedsNumberOfDocsInDirectoryFile(t *testing.T) {
	// Given 5 documents in directory file
	f := fixture()
	files := map[string]string {
		"1" : "first",
		"2" : "second",
		"3" : "third",
		"4" : "fourth",
		"5" : "fifth",
	}

	for name, content := range files {
		f.NewFile(name, content)
	}

	// When documents 0-6 are requested
	docs, err := f.Documents(0, 6)
	if err != nil {
		t.Fatalf("Unexpected err: ", err)
	}

	// Expect that 5 documents are returned
	if n := len(docs); n != 5 {
		t.Errorf("Expected 5; Got %d", n)
	}

	for i, doc := range docs {
		exp := files[strconv.Itoa(i+1)]
		act := string(doc.Text)
		if err := Compare(exp, act); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDocuments_WhenEmpty(t *testing.T) {
	// Given no documents in directory file
	f := fixture()

	// When documents 0-1 are requested
	docs, err := f.Documents(0, 6)
	if err != nil {
		t.Fatalf("Unexpected err: ", err)
	}

	// Expect that no documents are returned
	if n := len(docs); n != 0 {
		t.Errorf("Expected 0; Got %d", n)
	}
}
