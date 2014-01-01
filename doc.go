package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"strings"
)

type document struct {
	Path, Title, Text string
}

func validateDocument(lines []string) error {
	invalidMarkdownDocumentError := fmt.Errorf("Invalid markdown document")
	if len(lines) < 3 {
		return invalidMarkdownDocumentError
	}

	// make sure the second line is `==========` or some such
	for _, char := range lines[1] {
		if char != '=' {
			return invalidMarkdownDocumentError
		}
	}
	return nil
}

func fromLines(lines []string) string {
	out := ""
	for _, line := range lines {
		out += line + "\n"
	}
	return out
}

func parseDoc(path string, r io.Reader) (*document, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	if err := validateDocument(lines); err != nil {
		return nil, err
	}

	return &document{
		Path:  path,
		Title: lines[0],
		Text:  string(blackfriday.MarkdownCommon([]byte(fromLines(lines[2:])))),
	}, nil
}
