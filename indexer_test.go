package tinysearch

import (
	"reflect"
	"strings"
	"testing"
)

func TestUpdatePostingsList(t *testing.T) {

	// given
	setup()
	collection := []string{
		"Do you quarrel, sir?",
		"Quarrel sir! no, sir!",
		"If you do, sir, I am for you: I serve as good a man as you.",
		"No better.",
		"Well, sir",
	}

	indexer := NewIndexer(NewTokenizer())

	// when
	for i, doc := range collection {
		indexer.update(docID(i), strings.NewReader(doc))
	}

	// then
	dictionary := map[string]PostingsList{
		"a":       NewPostingsList(NewPosting(2, 12)),
		"am":      NewPostingsList(NewPosting(2, 5)),
		"as":      NewPostingsList(NewPosting(2, 10, 14)),
		"better":  NewPostingsList(NewPosting(3, 1)),
		"do":      NewPostingsList(NewPosting(0, 0), NewPosting(2, 2)),
		"for":     NewPostingsList(NewPosting(2, 6)),
		"good":    NewPostingsList(NewPosting(2, 11)),
		"i":       NewPostingsList(NewPosting(2, 4, 8)),
		"if":      NewPostingsList(NewPosting(2, 0)),
		"man":     NewPostingsList(NewPosting(2, 13)),
		"no":      NewPostingsList(NewPosting(1, 2), NewPosting(3, 0)),
		"quarrel": NewPostingsList(NewPosting(0, 2), NewPosting(1, 0)),
		"serve":   NewPostingsList(NewPosting(2, 9)),
		"sir":     NewPostingsList(NewPosting(0, 3), NewPosting(1, 1, 3), NewPosting(2, 3), NewPosting(4, 1)),
		"well":    NewPostingsList(NewPosting(4, 0)),
		"you":     NewPostingsList(NewPosting(0, 1), NewPosting(2, 1, 7, 15)),
	}

	expected := &Index{
		Dictionary:     dictionary,
		DocsLength:     map[docID]int{0: 4, 1: 4, 2: 16, 3: 2, 4: 2},
		TotalDocsCount: 5,
	}

	if !reflect.DeepEqual(indexer.index, expected) {
		t.Errorf("wrong index. \n\nexpected: \n%v\n\n got:\n%v\n", expected, indexer.index)
	}
}
