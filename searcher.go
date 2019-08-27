package tinysearch

import (
	"fmt"
	"math"
	"sort"
)

type ScoreDoc struct {
	docID DocumentID
	score float64
}

func (d ScoreDoc) String() string {
	return fmt.Sprintf("docId: %v, score: %v", d.docID, d.score)
}

type ScoreDocs map[DocumentID]*ScoreDoc

func NewScoreDocs(size int) ScoreDocs {
	return make(ScoreDocs, size)
}

// Add add and update score.
func (d ScoreDocs) Add(docID DocumentID, score float64) {
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

// 検索に使用するポスティングリストのポインタを取得する
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

// 検索を実行し、マッチしたドキュメントをスコア付きで返す
func (s *Searcher) search(query []string) ScoreDocs {

	// クエリに含まれる単語のポスティングリストが
	// ひとつも存在しない場合、0件で終了する
	if s.openCursors(query) == 0 {
		return ScoreDocs{}
	}

	c := s.cursors[0] // 一番短いカーソルを選択
	cursors := s.cursors[1:]

	// 結果を格納する構造体の初期化
	docs := NewScoreDocs(s.index.TotalDocsCount)

	for !c.Empty() {

		var nextDocId DocumentID

		for _, cursor := range cursors {
			// docId以上になるまで読み進める
			if cursor.NextDoc(c.DocId()); cursor.Empty() {
				return docs
			}
			// docIdが一致しなければ
			if cursor.DocId() != c.DocId() {
				nextDocId = cursor.DocId()
				break
			}
		}

		if nextDocId > 0 {
			// nextDocId以上になるまで読みすすめる
			if c.NextDoc(nextDocId); c.Empty() {
				return docs
			}
		} else {
			// 結果を格納
			docs.Add(c.DocId(), s.calcScore())
			c.Next()
		}
	}

	return docs
}

// 検索を実行し、スコアが高い順にK件結果を返す
func (s *Searcher) searchTopK(query []string, k int) *TopDocs {

	docs := s.search(query)

	// 結果をsliceに変換しスコアの降順でソートする
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
		totalHits: len(docs),
		scoreDocs: results,
	}
}
