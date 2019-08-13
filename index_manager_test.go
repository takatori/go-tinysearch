package tinysearch

import (
	"reflect"
	"strings"
	"testing"
)


func TestTextToWordSequence(t *testing.T) {

	document := "Do you quarrel, sir? Quarrel sir! no, sir! " +
		"If you do, sir, I am for you: " +
		"I serve as good a man as you. No better. Well, sir"

	expected := []string{
		"do", "you", "quarrel", "sir", "quarrel", "sir", "no", "sir",
		"if", "you", "do", "sir", "i", "am", "for", "you",
		"i", "serve", "as", "good", "a", "man", "as",
		"you", "no", "better", "well", "sir"}

	actual := TextToWordSequence(document)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nexpected: %v\n  actual: %v", expected, actual)
	}
}

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

	im := NewIndexManager()

	// when
	for i, doc := range collection {
		if err := im.updatePostingsList(int64(i+1), strings.NewReader(doc)); err != nil {
			t.Fatalf("failed to create index. beause: %v", err)
		}
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
		dictionary:     dictionary,
		documentLength: map[int64]int{1: 4, 2: 4, 3: 16, 4: 2, 5: 2},
		documentCount:  5,
	}

	if !reflect.DeepEqual(im.index, expected) {
		t.Errorf("wrong index. \n\nexpected: \n%v\n\n got:\n%v\n", expected, im.index)
	}
}
