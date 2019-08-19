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

	im := NewIndexer(NewTokenizer())

	// when
	for i, doc := range collection {
		im.update(docID(i+1), strings.NewReader(doc))
	}

	// then
	dictionary := map[string]PostingsList{
		"a":       NewPostingsList(NewPosting(3, []int{12})),
		"am":      NewPostingsList(NewPosting(3, []int{5})),
		"as":      NewPostingsList(NewPosting(3, []int{10, 14})),
		"better":  NewPostingsList(NewPosting(4, []int{1})),
		"do":      NewPostingsList(NewPosting(1, []int{0}), NewPosting(3, []int{2})),
		"for":     NewPostingsList(NewPosting(3, []int{6})),
		"good":    NewPostingsList(NewPosting(3, []int{11})),
		"i":       NewPostingsList(NewPosting(3, []int{4, 8})),
		"if":      NewPostingsList(NewPosting(3, []int{0})),
		"man":     NewPostingsList(NewPosting(3, []int{13})),
		"no":      NewPostingsList(NewPosting(2, []int{2}), NewPosting(4, []int{0})),
		"quarrel": NewPostingsList(NewPosting(1, []int{2}), NewPosting(2, []int{0})),
		"serve":   NewPostingsList(NewPosting(3, []int{9})),
		"sir":     NewPostingsList(NewPosting(1, []int{3}), NewPosting(2, []int{1, 3}), NewPosting(3, []int{3}), NewPosting(5, []int{1})),
		"well":    NewPostingsList(NewPosting(5, []int{0})),
		"you":     NewPostingsList(NewPosting(1, []int{1}), NewPosting(3, []int{1, 7, 15})),
	}

	expected := &Index{
		Dictionary: dictionary,
		DocLength:  map[docID]int{1: 4, 2: 4, 3: 16, 4: 2, 5: 2},
		DocCount:   5,
	}

	if !reflect.DeepEqual(im.index, expected) {
		t.Errorf("wrong index. \n\nexpected: \n%v\n\n got:\n%v\n", expected, im.index)
	}
}
