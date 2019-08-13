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
	dictionary map[string]PostingsList // 辞書
	docLength  map[documentID]int      // 各ドキュメントのサイズ
	docCount   int                     // ドキュメントの総数
}

// NewIndex create a new index.
func NewIndex() *Index {
	dict := make(map[string]PostingsList)
	length := make(map[documentID]int)
	return &Index{
		dictionary: dict,
		docLength:  length,
		docCount:   0,
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

	return fmt.Sprintf("docLength: %v, docCount: %v, dictionary: %v", idx.docLength, idx.docCount, strings.Join(str, "\n"))
}

// ポスティングリスト
type PostingsList struct {
	*list.List
}

func (pl PostingsList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
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

func (pl PostingsList) add(p *Posting) {
	pl.PushBack(p)
}

func (pl PostingsList) last() *Posting {
	e := pl.List.Back()
	if e == nil {
		return nil
	}
	return e.Value.(*Posting)
}

// ドキュメントをポスティングリストに追加
// ポスティングリストの最後を取得してドキュメントIDが
// 一致していなければ、ポスティングを追加
// 一致していれば、offsetを追加
func (pl PostingsList) Add(new *Posting) {
	last := pl.last()
	if last == nil || last.docID != new.docID {
		pl.add(new)
		return
	}
	last.offsets = append(last.offsets, new.offsets...)
	last.termFrequency++
}

// ドキュメントID
type documentID int64

// ポスティング
type Posting struct {
	docID         documentID // ドキュメントID
	offsets       []int      // 出現位置
	termFrequency int        // ドキュメント内の用語の出現回数
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p.docID, p.termFrequency, p.offsets)
}

// ポスティングを作成する
func NewPosting(docId documentID, offsets []int) *Posting {
	return &Posting{docId, offsets, len(offsets)}
}
