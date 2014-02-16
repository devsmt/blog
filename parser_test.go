package main

import (
	"testing"
	"strings"
	"fmt"
)

func expectStr(exp, act string) error {
	if exp != act {
		return fmt.Errorf("Expected '%s'; Got '%s'", exp, act)
	}
	return nil
}

func testParse(input string, metadata map[string]string, body string) error {
	p := &MetadataParser{}
	r := strings.NewReader(input)
	doc, err := p.Parse(r)

	if err != nil {
		return fmt.Errorf("Unexpected error: %v", err)
	}

	for k, v := range metadata {
		if err := expectStr(v, doc.Metadata[k]); err != nil {
			return err
		}
	}

	if err := expectStr(body, doc.String()); err != nil {
		return err
	}

	return nil
}

func TestParse_NoMetadata(t *testing.T) {
	if err := testParse(
		"The only line",
		map[string]string{},
		"The only line",
	); err != nil {
		t.Fatal(err)
	}
}

func TestParse_MultipleLinesInBody(t *testing.T) {
	if err := testParse(
		"line 1\nline 2",
		map[string]string{},
		"line 1\nline 2",
	); err != nil {
		t.Fatal(err)
	}
}

func TestParse_WithMetadata(t *testing.T) {
	if err := testParse(
		"Title:SomeTitle\nAuthor:SomeName\nThis is the body",
		map[string]string {
			"Title" : "SomeTitle",
			"Author": "SomeName",
		},
		"This is the body",
	); err != nil {
		t.Fatal(err)
	}
}

func TestParse_MultipleWordsInValue(t *testing.T) {
	if err := testParse(`Title:Some Title
This is the body`, map[string]string {
			"Title" : "Some Title",
		}, "This is the body"); err != nil {
		t.Fatal(err)
	}
}

// Expect error? What should happen in this scenario
func TestParse_MultipleWordsInKey(t *testing.T) {

}

func TestParse_SpacesAroundDelimeter(t *testing.T) {
	if err := testParse(`Title : Some Title
This is the body`, map[string]string {
			"Title" : "Some Title",
		}, "This is the body"); err != nil {
		t.Fatal(err)
	}
}

func TestParse_TrailingValueSpaceAndLeadingKeySpace(t *testing.T) {
	if err := testParse(" Title : Some Title \nThis is the body",
		map[string]string {
			"Title" : "Some Title",
		}, "This is the body"); err != nil {
		t.Fatal(err)
	}
}
