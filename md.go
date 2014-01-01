package main

import (
	"github.com/russross/blackfriday"
	"io"
	"text/template"
)

type MarkdownRenderer struct {
	template *template.Template
}

func NewMarkdownRenderer(template *template.Template) *MarkdownRenderer {
	return &MarkdownRenderer{template: template}
}

func (r *MarkdownRenderer) Write(w io.Writer, input []byte) error {
	input = blackfriday.MarkdownCommon(input)
	return r.template.Execute(w, string(input))
}

func (r *MarkdownRenderer) WriteHome(w io.Writer, docs []*document) error {
	for _, doc := range docs {
		doc.Body = blackfriday.MarkdownCommon(doc.Body)
	}
	return r.template.Execute(w, docs)
}
