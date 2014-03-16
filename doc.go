package main

import (
	"github.com/russross/blackfriday"
	"strings"
)

type Document struct {
	Metadata map[string]string
	Text     []byte
}

func (d *Document) String() string { return string(d.Text) }
func (d *Document) Snippet() string {
	if end := strings.Index(string(d.Text), "<!-- more -->"); end != -1 {
		return markdown(d.Text[0:end])
	}
	return ""
}

func (d *Document) Html() string {
	return markdown(d.Text)
}

func markdown(data []byte) string {
	return string(blackfriday.MarkdownCommon(data))
}
