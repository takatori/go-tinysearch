package tinysearch

import (
	"math"
	"sort"
)

type ScoreDoc struct {
	docID docID
	score float64
}

type ScoreDocs map[docID]*ScoreDoc

func NewScoreDocs(size int) ScoreDocs {
	return make(ScoreDocs, size)
}

// Add add and update score.
func (results ScoreDocs) Add(docID docID, score float64) {
	if _, ok := results[docID]; !ok {
		results[docID] = &ScoreDoc{
			docID: docID,
			score: score,
		}
	}
	// すでに結果に存在する場合はスコアを更新
	results[docID].score += score
}

type TopDocs struct {
	totalHits int
	scoreDocs []*ScoreDoc
}

// tfの計算
func calcTF(termCount int) float64 {
	if termCount <= 0 {
		return 0
	}
	return math.Log2(float64(termCount)) + 1
}

// idfの計算
func calIDF(N, df int) float64 {
	return math.Log2(float64(N) / float64(df))
}

// tf-idfスコアを計算する
// 全用語で検索が必要になるため遅い
func calcScore(idx *Index, query []string) ScoreDocs {

	results := NewScoreDocs(idx.TotalDocsCount)
	// クエリに含まれる用語ごとにスコアを計算
	for _, term := range query {
		if postingsList, ok := idx.Dictionary[term]; ok {
			idf := calIDF(idx.TotalDocsCount, postingsList.Len())
			for c := postingsList.newCursor(); c.nonNil(); c = c.next() {
				posting := c.Posting()
				docID := posting.DocID
				score := calcTF(posting.TermFrequency) * idf
				results.Add(docID, score)
			}
		}
	}

	return results
}

// 検索を実行する
func search(idx *Index, query []string, k int) TopDocs {

	// TODO: フィルタリングとスコア付を分離する
	// calcScoreはスコアの計算のみを行うように修正する
	docs := calcScore(idx, query)

	// 結果をarray型に変換後スコアの降順でソートして返す
	list := make([]*ScoreDoc, 0, len(docs))
	for _, v := range docs {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].score > list[j].score
	})

	return TopDocs{
		totalHits: len(list),
		scoreDocs: list[:k],
	}
}
