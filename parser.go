package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type MetadataParser struct {
}

type metadata struct {
	k, v string
}

func parseMetadata(line string) *metadata {
	index := strings.Index(line, ":")
	if index == -1 {
		return nil
	}

	return &metadata{
		k: strings.TrimFunc(line[:index], unicode.IsSpace),
		v: strings.TrimFunc(line[index+1:], unicode.IsSpace),
	}
}

func (p *MetadataParser) Parse(r io.Reader) (*Document, error) {
	s := bufio.NewScanner(r)
	d := &Document{ Metadata: make(map[string]string) }
	for s.Scan() {
		meta := parseMetadata(s.Text())
		if meta == nil { // TODO if we get here, this line should be added as text
			break
		}
		d.Metadata[meta.k] = meta.v
	}
	for {
		d.Text = append(d.Text, []byte(s.Text()+"\n")...)
		if !s.Scan() {
			break
		}
	}
	d.Text = d.Text[:len(d.Text)-1] // remove trailing new line
	return d, nil
}
