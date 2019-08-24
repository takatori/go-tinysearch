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
	Dictionary     map[string]PostingsList // 辞書
	TotalDocsCount int                     // ドキュメントの総数
	DocsLength     map[docID]int           // 各ドキュメントのサイズ
}

// NewIndex create a new index.
func NewIndex() *Index {
	dict := make(map[string]PostingsList)
	length := make(map[docID]int)
	return &Index{
		Dictionary:     dict,
		TotalDocsCount: 0,
		DocsLength:     length,
	}
}

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

	return fmt.Sprintf("DocsLength: %v\nTotalDocsCount:%v\nDictionary:\n%v\n",
		idx.DocsLength, idx.TotalDocsCount, strings.Join(str, "\n"))
}
