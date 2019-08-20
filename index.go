package tinysearch

import (
	"fmt"
	"sort"
	"strings"
)

// Index represents a inverted index.
// 転地インデックス
// 注意:mapを使用しているのでマルチスレッドには対応していない
type Index struct {
	Dictionary map[string]PostingsList `json:"Dictionary"` // 辞書
	DocCount   int                     `json:"DocCount"`   // ドキュメントの総数
	DocLength  map[docID]int           `json:"DocLength"`  // 各ドキュメントのサイズ
}

// NewIndex create a new index.
func NewIndex() *Index {
	dict := make(map[string]PostingsList)
	length := make(map[docID]int)
	return &Index{
		Dictionary: dict,
		DocCount:   0,
		DocLength:  length,
	}
}

/*
func (idx Index) UnmarshalJSON(b []byte) error {
	i := NewIndex()
	err := json.Unmarshal(b, i)
	fmt.Println(i)
	return err
}*/

func (idx Index) String() string {

	keys := make([]string, 0, len(idx.Dictionary))

	for k := range idx.Dictionary {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	str := make([]string, len(keys))

	for i, k := range keys {
		if postingList, ok := idx.Dictionary[k]; ok {
			str[i] = fmt.Sprintf("'%s'->%s", k, postingList.String())
		}
	}

	return fmt.Sprintf("DocLength: %v, DocCount: %v, Dictionary: %v",
		idx.DocLength, idx.DocCount, strings.Join(str, "\n"))
}
