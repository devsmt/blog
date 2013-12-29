package dir

import (
	"io"
	"time"
	"encoding/json"
	"log"
)

type Entry struct {
	Path string
	Time time.Time
}

func (e *Entry) String() string {
	return e.Path + "; " + e.Time.String()
}

type Directory interface {
	Entries() []*Entry
	// Reloads the directory
	Refresh(dirfile io.Reader)
}

type directory struct {
	entries []*Entry
}

func (d *directory) Entries() []*Entry {
	return d.entries
}

func (d *directory) Refresh(src io.Reader) {
	temp := []*Entry{}
	if err := json.NewDecoder(src).Decode(&temp); err != nil {
		log.Println("Error decoding dirfile:", err)
		return // don't update d.entries if there was an error parsing dirfile
	}
	d.entries = temp
}

func NewDirectory() Directory {
	return &directory{entries: []*Entry{}}
}
