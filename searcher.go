package tinysearch

import (
	"fmt"
	"math"
	"sort"
)

type ScoreDoc struct {
	docID docID
	score float64
}

func (d ScoreDoc) String() string {
	return fmt.Sprintf("docId: %v, score: %v", d.docID, d.score)
}

type ScoreDocs map[docID]*ScoreDoc

func NewScoreDocs(size int) ScoreDocs {
	return make(ScoreDocs, size)
}

// Add add and update score.
func (d ScoreDocs) Add(docID docID, score float64) {
	if _, ok := d[docID]; !ok {
		d[docID] = &ScoreDoc{
			docID: docID,
			score: score,
		}
	}
	// すでに結果に存在する場合はスコアを更新
	d[docID].score += score
}

type TopDocs struct {
	totalHits int
	scoreDocs []*ScoreDoc
}

func (t *TopDocs) String() string {
	return fmt.Sprintf("\ntotal hits: %v\nresults: %v\n",
		t.totalHits, t.scoreDocs)
}

type Searcher struct {
	index         *Index
	postingsLists []PostingsList
	cursors       []*Cursor
}

func NewSearcher(idx *Index) *Searcher {
	return &Searcher{index: idx}
}

// TFの計算
func calcTF(termCount int) float64 {
	if termCount <= 0 {
		return 0
	}
	return math.Log2(float64(termCount)) + 1
}

// IDFの計算
func calIDF(N, df int) float64 {
	return math.Log2(float64(N) / float64(df))
}

// tf-idfスコアを計算する
func (s *Searcher) calcScore() float64 {
	var score float64
	for i := 0; i < len(s.cursors); i++ {
		score += calcTF(s.cursors[i].Posting().TermFrequency) *
			calIDF(s.index.TotalDocsCount, s.postingsLists[i].Len())
	}
	return score
}

// カーソルを開く
func (s *Searcher) openCursors(query []string) int {
	// ポスティングリストの取得
	postingLists := make([]PostingsList, 0, len(query))
	for _, term := range query {
		if postingList, ok := s.index.Dictionary[term]; ok {
			postingLists = append(postingLists, postingList)
		}
	}
	if len(postingLists) == 0 {
		return 0
	}
	// ポスティングリストの短い順にソート
	sort.Slice(postingLists, func(i, j int) bool {
		return postingLists[i].Len() < postingLists[j].Len()
	})
	// 各ポスティングリストに対するcursorの取得
	cursors := make([]*Cursor, len(postingLists))
	for i, postingList := range postingLists {
		cursors[i] = postingList.openCursor()
	}
	s.postingsLists = postingLists
	s.cursors = cursors
	return len(cursors)
}

func (s *Searcher) search(query []string) ScoreDocs {

	// クエリに含まれる単語のポスティングリストが
	// ひとつも存在しない場合、0件で終了する
	if s.openCursors(query) == 0 {
		return ScoreDocs{}
	}

	c0 := s.cursors[0] // TODO: rename
	cursors := s.cursors[1:]

	// 結果を格納する構造体の初期化
	docs := NewScoreDocs(s.index.TotalDocsCount)

	for !c0.Empty() {

		var nextDocId docID

		for _, cursor := range cursors {
			// docId以上になるまで読み進める
			if cursor.NextDoc(c0.DocId()); cursor.Empty() {
				return docs
			}
			// docIdが一致しなければ
			if cursor.DocId() != c0.DocId() {
				nextDocId = cursor.DocId()
				break
			}
		}

		if nextDocId > 0 {
			// nextDocId以上になるまで読みすすめる
			if c0.NextDoc(nextDocId); c0.Empty() {
				return docs
			}
		} else {
			// 結果を格納
			docs.Add(c0.DocId(), s.calcScore())
			c0.Next()
		}
	}

	return docs
}

// 検索を実行する
func (s *Searcher) searchTopK(query []string, k int) *TopDocs {

	docs := s.search(query)

	// 結果をarray型に変換後スコアの降順でソートして返す
	results := make([]*ScoreDoc, 0, len(docs))
	for _, v := range docs {
		results = append(results, v)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// 上位k件のみ取得
	if len(results) > k {
		results = results[:k]
	}

	return &TopDocs{
		totalHits: len(results),
		scoreDocs: results,
	}
}
