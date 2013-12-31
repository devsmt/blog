package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io"
)

type MarkdownRenderer struct {
        Header []byte
        Footer []byte
}

func (r *MarkdownRenderer) Render(input []byte) []byte {
        return append(append(r.Header, blackfriday.MarkdownCommon(input)...), r.Footer...)
}

func (r *MarkdownRenderer) RenderHome(docs []*document) []byte {
        out := ""
        for _, doc := range docs {
                out += fmt.Sprintf(`<h2><a href="/%s">%s</a></h2>`, doc.Path, doc.Title)
                out += string(doc.Body)
        }
        return r.Render([]byte(out))
}

func (r *MarkdownRenderer) Write(w io.Writer, input []byte) error {
        _, err := w.Write(r.Render(input))
        return err
}

func (r *MarkdownRenderer) WriteHome(w io.Writer, docs []*document) error {
        _, err := w.Write(r.RenderHome(docs))
        return err
}

var defaultRenderer = &MarkdownRenderer{Header:[]byte{}, Footer:[]byte{}}

func SetHeader(header []byte) {
	defaultRenderer.Header = header
}

func SetFooter(footer []byte) {
	defaultRenderer.Footer = footer
}

func Write(w io.Writer, input []byte) error {
	_, err := w.Write(defaultRenderer.Render(input))
	return err
}

func WriteHome(w io.Writer, docs []*document) error {
	_, err := w.Write(defaultRenderer.RenderHome(docs))
	return err
}
