package main

import "github.com/russross/blackfriday"

type Document struct {
	Metadata map[string]string
	Text     []byte
}

func (d *Document) String() string { return string(d.Text) }
func (d *Document) Html() string {
	return string(blackfriday.MarkdownCommon(d.Text))
}
