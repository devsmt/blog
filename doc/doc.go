package doc

import (
	"bufio"
	"github.com/russross/blackfriday"
	"io"
	"strings"
	"unicode"
)

type Document struct {
	Metadata      map[string]string
	text          string
	snippetLength int
	Path          string
}

func (d *Document) Snippet() string {
	if d.snippetLength >= 0 && d.snippetLength <= len(d.text) {
		return d.text[:d.snippetLength]
	}
	return d.text
}

func (d *Document) Text() string {
	return d.text
}

func metadata(line string) (k, v string) {
	if i := strings.Index(line, ":"); i != -1 {
		return line[:i], strings.TrimLeftFunc(line[i+1:], unicode.IsSpace)
	}
	return "", ""
}

func isMetadata(k, v string) bool {
	return k != "" && v != ""
}

func Parse(r io.Reader) (*Document, error) {
	d := new(Document)
	d.snippetLength = -1
	d.Metadata = map[string]string{}

	data := []byte{}
	continueParsingMetadata := true
	for s := bufio.NewScanner(r); s.Scan(); {
		if err := s.Err(); err != nil {
			return nil, err
		}

		if continueParsingMetadata {
			k, v := metadata(s.Text())
			if isMetadata(k, v) {
				d.Metadata[k] = v
				continue
			}
			continueParsingMetadata = false
		}

		data = append(data, s.Bytes()...)
		data = append(data, '\n')
	}

	if len(data) > 0 {
		d.text = string(blackfriday.MarkdownCommon(data))
		d.text = d.text[:len(d.text)]
		d.snippetLength = strings.Index(d.text, "<!-- more -->")
	}
	return d, nil
}
