package tinysearch

import (
	"testing"
)

func TestSearch(t *testing.T) {

	// given
	dictionary := map[string]PostingsList{
		"a":      NewPostingsList(NewPosting(2, 12)),
		"am":     NewPostingsList(NewPosting(2, 5)),
		"as":     NewPostingsList(NewPosting(2, 10, 14)),
		"better": NewPostingsList(NewPosting(3, 1)),
		"do": NewPostingsList(
			NewPosting(0, 0),
			NewPosting(2, 2)),
		"for":  NewPostingsList(NewPosting(2, 6)),
		"good": NewPostingsList(NewPosting(2, 11)),
		"i":    NewPostingsList(NewPosting(2, 4, 8)),
		"if":   NewPostingsList(NewPosting(2, 0)),
		"man":  NewPostingsList(NewPosting(2, 13)),
		"no": NewPostingsList(
			NewPosting(1, 2),
			NewPosting(3, 0)),
		"quarrel": NewPostingsList(
			NewPosting(0, 2), NewPosting(1, 0)),
		"serve": NewPostingsList(NewPosting(2, 9)),
		"sir": NewPostingsList(
			NewPosting(0, 3),
			NewPosting(1, 1, 3),
			NewPosting(2, 3),
			NewPosting(4, 1)),
		"well": NewPostingsList(
			NewPosting(4, 0)),
		"you": NewPostingsList(
			NewPosting(0, 1),
			NewPosting(2, 1, 7, 15)),
	}

	idx := &Index{
		Dictionary:     dictionary,
		TotalDocsCount: 5,
	}
	s := NewSearcher(idx)

	query := []string{"quarrel", "sir"}

	// when
	actual := s.searchTopK(query, 4)

	expected := []docID{1, 0}

	// TODO: 長さとスコアのチェック
	for i, docID := range expected {
		if actual.scoreDocs[i].docID != docID {
			t.Fatalf("\ngot:\n%v\nexpected: %v\n", actual, expected)
		}
	}
}
