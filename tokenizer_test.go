package tinysearch

import (
	"reflect"
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

	tokenizer := &DefaultTokenizer{}
	actual := tokenizer.TextToWordSequence(document)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nexpected: %v\n  actual: %v", expected, actual)
	}
}
