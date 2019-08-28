package tinysearch

import (
	"container/list"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// 転地インデックス
// 注意: mapを使用しているのでマルチスレッドには対応していない
type Index struct {
	Dictionary     map[string]PostingsList // 辞書
	TotalDocsCount int                     // ドキュメントの総数
}

// NewIndex create a new index.
func NewIndex() *Index {
	dict := make(map[string]PostingsList)
	return &Index{
		Dictionary:     dict,
		TotalDocsCount: 0,
	}
}

func (idx Index) String() string {
	var padding int
	keys := make([]string, 0, len(idx.Dictionary))
	for k := range idx.Dictionary {
		l := utf8.RuneCountInString(k)
		if padding < l {
			padding = l
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	strs := make([]string, len(keys))
	format := "  [%-" + strconv.Itoa(padding) + "s] -> %s"
	for i, k := range keys {
		if postingList, ok := idx.Dictionary[k]; ok {
			strs[i] = fmt.Sprintf(format, k, postingList.String())
		}
	}
	return fmt.Sprintf("total documents : %v\ndictionary:\n%v\n",
		idx.TotalDocsCount, strings.Join(strs, "\n"))
}

// ドキュメントID
type DocumentID int64

// ポスティング
type Posting struct {
	DocID         DocumentID // ドキュメントID
	Positions     []int      // 用語の出現位置
	TermFrequency int        // ドキュメント内の用語の出現回数
}

// ポスティングを作成する
func NewPosting(docID DocumentID, positions ...int) *Posting {
	return &Posting{docID, positions, len(positions)}
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v,%v,%v)",
		p.DocID, p.TermFrequency, p.Positions)
}

// ポスティングリスト
type PostingsList struct {
	*list.List
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

func (pl PostingsList) MarshalJSON() ([]byte, error) {

	postings := make([]*Posting, 0, pl.Len())

	for e := pl.Front(); e != nil; e = e.Next() {
		postings = append(postings, e.Value.(*Posting))
	}
	return json.Marshal(postings)
}

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

func (pl PostingsList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}

// 現在の読み込み位置を表すポインタ
type Cursor struct {
	current *list.Element
}

func (pl PostingsList) OpenCursor() *Cursor {
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
