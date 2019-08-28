package tinysearch

import (
	"reflect"
	"testing"
)

func TestSearchTopK(t *testing.T) {

	reader := NewIndexReader("testdata/index")
	s := NewSearcher(reader) // ❶ searcherの初期化

	actual := s.searchTopK([]string{"quarrel", "sir"}, 1) // ❷ 検索の実行

	expected := &TopDocs{2,
		[]*ScoreDoc{
			{2, 1.9657842846620868},
		},
	}

	for !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got:\n%v\nexpected:%v\n", actual, expected)
	}
}
