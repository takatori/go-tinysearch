package tinysearch

import (
	"fmt"
	"sort"
)

// 検索結果を格納する構造体
type SearchResult struct {
	docID docID
	score float64
	title string
}

// String print search result info
func (r *SearchResult) String() string {
	return fmt.Sprintf("{DocID: %v, score: %.2f, title: %v}", r.docID, r.score, r.title)
}

type SearchResults map[docID]*SearchResult

func NewSearchResults(size int) SearchResults {
	return make(SearchResults, size)
}

// Add add and update score.
func (results SearchResults) Add(docID docID, score float64) {

	if _, ok := results[docID]; !ok {
		results[docID] = &SearchResult{
			docID: docID,
			score: score,
		}
	}

	// すでに結果に存在する場合はスコアを更新
	results[docID].score += score
}

// 結果をソートする
func (results SearchResults) Rank() []*SearchResult {

	list := make([]*SearchResult, 0, len(results))

	for _, v := range results {
		list = append(list, v)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].score > list[j].score
	})

	return list
}
