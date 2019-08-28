package tinysearch

import (
	"reflect"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {

	setup()
	collection := []string{
		"Do you quarrel, sir?",
		"Quarrel sir! no, sir!",
		"No better.",
		"Well, sir",
	}

	indexer := NewIndexer(NewTokenizer()) // ❶ インデックス構築器の初期化

	for i, doc := range collection {
		// ❷ インデックスにドキュメントを追加
		indexer.update(DocumentID(i), strings.NewReader(doc))
	}

	actual := indexer.index
	expected := &Index{
		Dictionary: map[string]PostingsList{
			"better": NewPostingsList(
				NewPosting(2, 1)),
			"do": NewPostingsList(
				NewPosting(0, 0)),
			"no": NewPostingsList(
				NewPosting(1, 2),
				NewPosting(2, 0)),
			"quarrel": NewPostingsList(
				NewPosting(0, 2),
				NewPosting(1, 0)),
			"sir": NewPostingsList(NewPosting(0, 3),
				NewPosting(1, 1, 3),
				NewPosting(3, 1)),
			"well": NewPostingsList(
				NewPosting(3, 0)),
			"you": NewPostingsList(
				NewPosting(0, 1)),
		},
		TotalDocsCount: 4,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("wrong index. \n\nwant: \n%v\n\n got:\n%v\n",
			expected, actual)
	}
}
