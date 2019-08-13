package tinysearch

import (
	"math"
)

// クエリベクトルの重みを計算する
func calcQueryWeight(N, df int) float64 {
	// tf * idf = 1 * log(N/dft)
	// todo: クエリベクトルに同一単語が複数含まれている場合はtfは1ではなくなる？
	return math.Log2(float64(N) / float64(df))
}

// ドキュメントベクトルの重みを計算する
func calcDocumentWeight(tf int) float64 {
	if tf <= 0 {
		return 0
	}
	return math.Log2(float64(tf)) + 1
}

// ベクトル空間スコアを計算する
// 全用語で検索が必要になるため遅い
func cosineScore(idx *Index, query []string) []*SearchResult {

	results := NewSearchResults(idx.documentCount)

	// クエリに含まれる用語ごとにスコアを計算
	for _, term := range query {

		postingsList, ok := idx.dictionary[term]
		if !ok {
			continue // 辞書に存在しない場合はスキップ
		}

		// クエリベクトルの重み計算
		wtq := calcQueryWeight(idx.documentCount, postingsList.Len())

		for e := postingsList.Front(); e != nil; e = e.Next() {
			posting := e.Value.(*Posting) // todo: キャストしなくても良いようにする
			docId := posting.docId
			wtd := calcDocumentWeight(len(posting.offsets))
			score := wtd * wtq
			results.Add(docId, score)
		}
	}
	// normalize
	for _, result := range results {
		result.score = result.score / float64(idx.documentLength[result.docId])
	}

	return results.Rank()
}
