package doc

import (
	"bufio"
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

		// check the current line for <!-- more -->; if found, set snippetLength
		if i := strings.Index(s.Text(), "<!-- more -->"); i != -1 {
			d.snippetLength = len(data) + i
		}
		data = append(data, s.Bytes()...)
		data = append(data, '\n')
	}

	if len(data) > 0 {
		d.text = string(data[:len(data)-1])
	}
	return d, nil
}
