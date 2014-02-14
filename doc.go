package main

import (
	"io"
	"io/ioutil"
)

type Document struct {
	Text []byte
}
func (d *Document) String() string { return string(d.Text) }

type FakeDocParser struct{}
func (p *FakeDocParser) Parse(r io.Reader) (*Document, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Document{Text: data}, nil
}
