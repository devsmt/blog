package main

import (
	"bufio"
)

type document struct {
        Title, Path string
        Body []byte
}

func parseDoc(path string, r io.Reader) *document {
        s := bufio.NewScanner(r)
        s.Scan()
        d := new(document)
        d.Title = s.Text() // set Title = first line
        s.Scan() // skip the second line (the ======= below the title)
        d.Path = path
        for s.Scan() { // add other lines to the body
                d.Body = append(d.Body, s.Bytes()...)
        }
        return d
}
