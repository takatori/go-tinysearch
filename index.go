package tinysearch

import (
	"container/list"
	"fmt"
	"sort"
	"strings"
)

// Index represents a inverted index.
// 転地インデックス
// 注意:mapを使用しているのでマルチスレッドには対応していない
type Index struct {
	dictionary     map[string]PostingsList // 辞書
	documentLength map[int64]int           // 各ドキュメントのサイズ
	documentCount  int                     // ドキュメントの総数
}

// NewIndex create a new index.
func NewIndex() *Index {
	dict := make(map[string]PostingsList)
	length := make(map[int64]int)
	return &Index{
		dictionary:     dict,
		documentLength: length,
		documentCount:  0,
	}
}

func (idx Index) String() string {

	keys := make([]string, 0, len(idx.dictionary))

	for k := range idx.dictionary {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	str := make([]string, len(keys))

	for i, k := range keys {
		if postingList, ok := idx.dictionary[k]; ok {
			str[i] = fmt.Sprintf("'%s'->%s", k, postingList.String())
		}
	}

	return fmt.Sprintf("documentLength: %v, documentCount: %v, dictionary: %v", idx.documentLength, idx.documentCount, strings.Join(str, "\n"))
}

// ポスティングリスト
type PostingsList struct {
	*list.List
}

func (l PostingsList) String() string {
	str := make([]string, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}

// ポスティングリストを作成する処理
func NewPostingsList(postings ...*Posting) PostingsList {
	l := list.New()
	for _, posting := range postings {
		l.PushBack(posting)
	}
	return PostingsList{l}
}

// ドキュメントをポスティングリストに追加
func (l PostingsList) Add(posting *Posting) {

	e := l.List.Back()
	if e == nil {
		l.List.PushBack(posting)
		return
	}

	lastPosting := e.Value.(*Posting)

	// ポスティングリストの最後を取得してドキュメントIDが一致していればoffsetを追加
	// TODO: 一番最後以外にいる可能性
	if lastPosting.docId == posting.docId {
		lastPosting.offsets = append(lastPosting.offsets, posting.offsets...)
		lastPosting.termFrequency++
		return
	}

	l.List.PushBack(posting)
}

// ポスティング
type Posting struct {
	docId         int64 // ドキュメントID
	offsets       []int // 出現位置
	termFrequency int   // ドキュメント内の用語の出現回数
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p.docId, p.termFrequency, p.offsets)
}

// ポスティングを作成する
func NewPosting(docId int64, offsets []int) *Posting {
	return &Posting{docId, offsets, len(offsets)}
}
