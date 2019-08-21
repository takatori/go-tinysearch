package tinysearch

import (
	"bufio"
	"io"
)

type Indexer struct {
	index     *Index
	tokenizer *Tokenizer
}

func NewIndexer(tokenizer *Tokenizer) *Indexer {
	return &Indexer{
		index:     NewIndex(),
		tokenizer: tokenizer,
	}
}

// ドキュメントをインデックスに追加する処理
func (idxr *Indexer) update(docID docID, reader io.Reader) {

	scanner := bufio.NewScanner(reader)
	scanner.Split(idxr.tokenizer.SplitFunc)
	var position int

	for scanner.Scan() {
		term := scanner.Text()
		if postingsList, ok := idxr.index.Dictionary[term]; !ok {
			// termをキーとするポスティングリストが存在しない場合は新規作成
			idxr.index.Dictionary[term] = NewPostingsList(NewPosting(docID, []int{position}))
		} else {
			// ポスティングリストがすでに存在する場合は追加
			postingsList.Add(NewPosting(docID, []int{position}))
		}
		position++
	}

	idxr.index.DocsCount++
	idxr.index.DocsLength[docID] = position
}
