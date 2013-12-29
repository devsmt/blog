package dir

import (
	"time"
	"strings"
	"fmt"
	"testing"
	"encoding/json"
	"io"
	"log"
	"bytes"
)

func validDirfile(entries []*Entry) io.Reader {
	data, err := json.Marshal(entries)
	if err != nil {
		log.Fatal("Unexpected err:", err)
	}
	return bytes.NewBuffer(data)
}

func entryListEq(el1, el2 []*Entry) bool {
	if len(el1) != len(el2) {
		return false
	}

	for i:=0; i<len(el1); i++ {
		if el1[i].Path != el2[i].Path {
			return false
		}
	}
	return true
}

func ExpectEntryListEq(el1, el2 []*Entry) error {
	if !entryListEq(el1, el2) {
		return fmt.Errorf("Expected lists to be equal:\n1: %v\n2: %v", el1, el2)
	}
	return nil
}

// Given dirfile contains valid JSON
// When Refresh() called
// Then expect Entries() returns correct entry list
func TestRefresh_GivenValidJSON_ExpectSuccessfulDecoding(t *testing.T) {
	testEntries := []*Entry {
		&Entry{Path: "1.txt"},
		&Entry{Path: "2.txt"},
	}

	// Given dirfile contains valid JSON
	dirfile := validDirfile(testEntries)

	d := NewDirectory()

	d.Refresh(dirfile)

	if err := ExpectEntryListEq(testEntries, d.Entries()); err != nil {
		t.Fatal(err)
	}
}

func makeEntries(n int) []*Entry {
	entries := make([]*Entry, n)
	for i:=0; i<n; i++ {
		entries[i] = &Entry{Path: fmt.Sprintf("%d.txt", i), Time: time.Now()}
	}
	return entries
}

// Given some entries in Directory
// And dirfile contains invalid JSON
// When Refresh() called
// Then expect Entries() remains unchanged
func TestRefresh_GivenInvalidJSON_ExpectNoChangeToEntries(t *testing.T) {
	// Given some entries in Directory
	entries := makeEntries(3)
	d := &directory{entries: entries}

	// And dirfile contains invalid JSON
	invalidJSON := `{"Path":`
	dirfile := strings.NewReader(invalidJSON)

	// When Refresh() called
	d.Refresh(dirfile)

	// Then expect Entries() unchanged
	if err := ExpectEntryListEq(entries, d.Entries()); err != nil {
		t.Fatal(err)
	}
}
