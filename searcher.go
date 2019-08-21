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

	results := NewSearchResults(idx.DocsCount)

	// クエリに含まれる用語ごとにスコアを計算
	for _, term := range query {

		postingsList, ok := idx.Dictionary[term]
		if !ok {
			continue // 辞書に存在しない場合はスキップ
		}

		// クエリベクトルの重み計算
		wtq := calcQueryWeight(idx.DocsCount, postingsList.Len())

		for cursor := postingsList.NewCursor(); cursor.Element != nil; cursor = cursor.next() {
			posting := cursor.Posting()
			docID := posting.DocID
			wtd := calcDocumentWeight(len(posting.Positions))
			score := wtd * wtq
			results.Add(docID, score)
		}
	}
	// normalize
	for _, result := range results {
		result.score = result.score / float64(idx.DocsLength[result.docID])
	}

	return results.Rank() // TODO: topk
}
