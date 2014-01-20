package doc

import (
	"github.com/russross/blackfriday"
	"testing"
	"strings"
	"log"
	"fmt"
	"unicode"
)

func setupParseTest(input string) *Document {
	d, err := Parse(strings.NewReader(input))
	if err != nil {
		log.Fatal("Unexpected err: ", err)
	}
	return d
}

func md(s string) string {
	return string(blackfriday.MarkdownCommon([]byte(s)))
}

func expect(exp, act string) error {
	if act != exp {
		return fmt.Errorf("Expected:\n%s\nGot:\n%s\n", exp, act)
	}
	return nil
}

func expectText(d *Document, exp string) error {
	return expect(md(exp), d.Text())
}

func expectSnippet(d *Document, exp string) error {
	exp = strings.TrimFunc(md(exp), unicode.IsSpace)
	act := strings.TrimFunc(d.Snippet(), unicode.IsSpace)
	return expect(exp, act)
}

func expectMetadata(d *Document, key, exp string) error {
	return expect(exp, d.Metadata[key])
}

func TestParse(t *testing.T) {
	input := "This is some text"
	d := setupParseTest(input)

	if err := expectText(d, input); err != nil {
		t.Fatal(err)
	}
}

func TestParseWithMultiLineInput(t *testing.T) {
	input := "line1\nline2"
	d := setupParseTest(input)

	if err := expectText(d, input); err != nil {
		t.Fatal(err)
	}
}

func TestSnippet(t *testing.T) {
	input := "This is some text.\n\n<!-- more -->\n\nThis is some more text"
	d := setupParseTest(input)

	if err := expectText(d, input); err != nil {
		t.Fatal(err)
	}
	if err := expectSnippet(d, "This is some text."); err != nil {
		t.Fatal(err)
	}
}

// test Snippet() when <!-- more --> doesn't exist;
// expect d.Snippet() == d.Text()
func TestSnippet_WithoutMoreTag(t *testing.T) {
	input := "This is some text"
	d := setupParseTest(input)

	if err := expectSnippet(d, input); err != nil {
		t.Fatal(err)
	}
}

func TestMetadata(t *testing.T) {
	input := "Title: This is the title\nThis is just some text"

	d := setupParseTest(input)
	if err := expectText(d, "This is just some text"); err != nil {
		t.Fatal(err)
	}

	if err := expectMetadata(d, "Title", "This is the title"); err != nil {
		t.Fatal(err)
	}
}

func TestMetadata_WithDifferentMetadata(t *testing.T) {
	input := "Title: a different title\nThis is just some text"

	d := setupParseTest(input)
	if err := expectText(d, "This is just some text"); err != nil {
		t.Fatal(err)
	}

	if err := expectMetadata(d, "Title", "a different title"); err != nil {
		t.Fatal(err)
	}
}

func TestMetadata_WithMultipleFields(t *testing.T) {
	input := "Title: title\nDate: date\nSome text"

	d := setupParseTest(input)
	if err := expectText(d, "Some text"); err != nil {
		t.Fatal(err)
	}

	if err := expectMetadata(d, "Title", "title"); err != nil {
		t.Fatal(err)
	}

	if err := expectMetadata(d, "Date", "date"); err != nil {
		t.Fatal(err)
	}
}

func TestMetadata_NoSpacesAroundDelimeter(t *testing.T) {
	input := "Title:title\nasdf"

	d := setupParseTest(input)
	if err := expectText(d, "asdf"); err != nil {
		t.Fatal(err)
	}

	if err := expectMetadata(d, "Title", "title"); err != nil {
		t.Fatal(err)
	}
}

// An empty body shouldn't crash anything.
func TestMetadata_NoBodyText(t *testing.T) {
	input := "Title: title"
	d := setupParseTest(input)
	if err := expectText(d, ""); err != nil {
		t.Fatal(err)
	}

	if err := expectMetadata(d, "Title", "title"); err != nil {
		t.Fatal(err)
	}
}
