package tinysearch

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSearchTopK(t *testing.T) {

	idxStr := `
{
    "TotalDocsCount":5,
    "Dictionary":{
	"a":      [{"DocID":3,"Positions":[12],"TermFrequency":1}],
	"am":     [{"DocID":3,"Positions":[5],"TermFrequency":1}],
	"as":     [{"DocID":3,"Positions":[10,14],"TermFrequency":2}],
	"better": [{"DocID":4,"Positions":[1],"TermFrequency":1}],
	"do":     [{"DocID":1,"Positions":[0],"TermFrequency":1},
               {"DocID":3,"Positions":[2],"TermFrequency":1}],
	"for":    [{"DocID":3,"Positions":[6],"TermFrequency":1}],
	"good":   [{"DocID":3,"Positions":[11],"TermFrequency":1}],
	"i":      [{"DocID":3,"Positions":[4,8],"TermFrequency":2}],
	"if":     [{"DocID":3,"Positions":[0],"TermFrequency":1}],
	"man":    [{"DocID":3,"Positions":[13],"TermFrequency":1}],
	"no":     [{"DocID":2,"Positions":[2],"TermFrequency":1},
               {"DocID":4,"Positions":[0],"TermFrequency":1}],
	"quarrel":[{"DocID":1,"Positions":[2],"TermFrequency":1},
               {"DocID":2,"Positions":[0],"TermFrequency":1}],
	"serve":  [{"DocID":3,"Positions":[9],"TermFrequency":1}],
	"sir":    [{"DocID":1,"Positions":[3],"TermFrequency":1},
               {"DocID":2,"Positions":[1,3],"TermFrequency":2},
               {"DocID":3,"Positions":[3],"TermFrequency":1},
               {"DocID":5,"Positions":[1],"TermFrequency":1}],
	"well":   [{"DocID":5,"Positions":[0],"TermFrequency":1}],
	"you":    [{"DocID":1,"Positions":[1],"TermFrequency":1},
               {"DocID":3,"Positions":[1,7,15],"TermFrequency":3}]
    }
}
`
	idx := NewIndex()
	if err := json.Unmarshal([]byte(idxStr), idx); err != nil {
		t.Fatalf("failed to unmarshal idxStr: %v", err)
	}
	s := NewSearcher(idx) // ❶ searcherの初期化

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
