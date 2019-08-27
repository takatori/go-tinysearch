package tinysearch

import (
	"container/list"
	"encoding/json"
	"fmt"
	"strings"
)

// ドキュメントID
type DocumentID int64

// ポスティング
type Posting struct {
	DocID         DocumentID // ドキュメントID
	Positions     []int      // 用語の出現位置
	TermFrequency int        // ドキュメント内の用語の出現回数
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p.DocID, p.TermFrequency, p.Positions)
}

// ポスティングを作成する
func NewPosting(docID DocumentID, positions ...int) *Posting {
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

// 現在の読み込み位置を表すポインタ
type Cursor struct {
	current *list.Element
}

func (pl PostingsList) openCursor() *Cursor {
	return &Cursor{pl.Front()}
}

func (c *Cursor) Next() {
	c.current = c.current.Next()
}

func (c *Cursor) NextDoc(id DocumentID) {
	for !c.Empty() && c.DocId() < id {
		c.Next()
	}
}

func (c *Cursor) Empty() bool {
	if c.current == nil {
		return true
	}
	return false
}

func (c *Cursor) Posting() *Posting {
	return c.current.Value.(*Posting)
}

func (c *Cursor) DocId() DocumentID {
	return c.current.Value.(*Posting).DocID
}

func (c *Cursor) String() string {
	return fmt.Sprint(c.Posting())
}
