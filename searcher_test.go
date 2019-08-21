package tinysearch

import (
	"testing"
)

func TestCosineScore(t *testing.T) {

	// given
	dictionary := map[string]PostingsList{
		"a":       NewPostingsList(NewPosting(3, []int{13})),
		"am":      NewPostingsList(NewPosting(3, []int{6})),
		"as":      NewPostingsList(NewPosting(3, []int{11, 15})),
		"better":  NewPostingsList(NewPosting(4, []int{2})),
		"do":      NewPostingsList(NewPosting(1, []int{1}), NewPosting(3, []int{3})),
		"for":     NewPostingsList(NewPosting(3, []int{7})),
		"good":    NewPostingsList(NewPosting(3, []int{12})),
		"i":       NewPostingsList(NewPosting(3, []int{5, 9})),
		"if":      NewPostingsList(NewPosting(3, []int{1})),
		"man":     NewPostingsList(NewPosting(3, []int{14})),
		"no":      NewPostingsList(NewPosting(2, []int{3}), NewPosting(4, []int{1})),
		"quarrel": NewPostingsList(NewPosting(1, []int{3}), NewPosting(2, []int{1})),
		"serve":   NewPostingsList(NewPosting(3, []int{10})),
		"sir":     NewPostingsList(NewPosting(1, []int{4}), NewPosting(2, []int{2, 4}), NewPosting(3, []int{4}), NewPosting(5, []int{2})),
		"well":    NewPostingsList(NewPosting(5, []int{1})),
		"you":     NewPostingsList(NewPosting(1, []int{2}), NewPosting(3, []int{2, 8, 16})),
	}

	idx := &Index{
		Dictionary: dictionary,
		DocsLength: map[docID]int{1: 4, 2: 4, 3: 16, 4: 2, 5: 2},
		DocsCount:  5,
	}

	query := []string{"quarrel", "sir"}

	// when
	actual := cosineScore(idx, query)

	expected := []docID{2, 1, 5, 3}

	// TODO: 長さとスコアのチェック
	for i, docID := range expected {
		if actual[i].docID != docID {
			t.Fatalf("\ngot:\n%v\nexpected: %v\n", actual, expected)
		}
	}
}
