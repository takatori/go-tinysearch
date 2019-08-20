package tinysearch

import (
	"container/list"
	"encoding/json"
	"fmt"
	"strings"
)

// ドキュメントID
type docID int64

// ポスティング
type Posting struct {
	DocID         docID `json:"DocID"`         // ドキュメントID
	Positions     []int `json:"Positions"`     // 用語の出現位置
	TermFrequency int   `json:"TermFrequency"` // ドキュメント内の用語の出現回数
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p.DocID, p.TermFrequency, p.Positions)
}

// ポスティングを作成する
func NewPosting(docID docID, positions []int) *Posting {
	return &Posting{docID, positions, len(positions)}
}

// ポスティングリスト
type PostingsList struct {
	*list.List
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

func (pl PostingsList) MarshalJSON() ([]byte, error) {

	postings := make([]*Posting, 0, pl.Len())

	for e := pl.Front(); e != nil; e = e.Next() {
		postings = append(postings, e.Value.(*Posting))
	}
	return json.Marshal(postings)
}

func (pl PostingsList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}

// ドキュメントをポスティングリストに追加
// ポスティングリストの最後を取得してドキュメントIDが
// 一致していなければ、ポスティングを追加
// 一致していれば、positionを追加
func (pl PostingsList) Add(new *Posting) {
	last := pl.last()
	if last == nil || last.DocID != new.DocID {
		pl.add(new)
		return
	}
	last.Positions = append(last.Positions, new.Positions...)
	last.TermFrequency++
}

// ポスティングリストを作成する処理
func NewPostingsList(postings ...*Posting) PostingsList {
	l := list.New()
	for _, posting := range postings {
		l.PushBack(posting)
	}
	return PostingsList{l}
}

// [For Search]
func (pl *PostingsList) UnmarshalJSON(b []byte) error {

	var postings []*Posting
	if err := json.Unmarshal(b, &postings); err != nil {
		return err
	}
	pl.List = list.New()
	for _, posting := range postings {
		pl.add(posting)
	}

	return nil
}

type Cursor struct {
	*list.Element
}

func (pl PostingsList) NewCursor() Cursor {
	return Cursor{pl.Front()}
}

func (c Cursor) next() Cursor {
	return Cursor{c.Next()}
}

func (c Cursor) Posting() *Posting {
	return c.Value.(*Posting)
}
