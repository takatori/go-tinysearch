package tinysearch

import (
	"fmt"
	"sort"
)

// 検索結果を格納する構造体
type SearchResult struct {
	docId int64
	score float64
	title string
}

// String print search result info

func (r *SearchResult) String() string {
	return fmt.Sprintf("{docId: %v, score: %.2f, title: %v}", r.docId, r.score, r.title)
}

type SearchResults map[int64]*SearchResult

func NewSearchResults(size int) SearchResults {
	return make(SearchResults, size)
}

// Add add and update score.
func (results SearchResults) Add(docId int64, score float64) {

	if _, ok := results[docId]; !ok {
		results[docId] = &SearchResult{
			docId: docId,
			score: score,
		}
	}

	// すでに結果に存在する場合はスコアを更新
	results[docId].score += score
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
