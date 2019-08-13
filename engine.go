package tinysearch

import (
	"database/sql"
	"io"
)

type Tokenizer interface {
	TextToWordSequence(string) []string
	SplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error)
}

// 検索エンジン
type Engine struct {
	tokenizer       Tokenizer
	indexer         *Indexer
	documentManager *DocumentManager
}

// NewSearchEngine(db) create a search engine.
// 検索エンジンを作成する
func NewSearchEngine(db *sql.DB) *Engine {
	tokenizer := &DefaultTokenizer{}
	return &Engine{
		tokenizer,
		NewIndexer(tokenizer),
		NewDocumentManager(db),
	}
}

// インデックスにドキュメントを追加する
func (e *Engine) AddDocument(title string, reader io.Reader) error {

	id, err := e.documentManager.save(title)
	if err != nil {
		return err
	}

	e.indexer.update(id, reader)
	return nil
}

// TODO
// インデックスをファイルに書き出す
// func (e *Engine) Commit() error {
//
// }

// 検索を実行する
func (e *Engine) Search(query string, k int) ([]*SearchResult, error) {

	terms := e.tokenizer.TextToWordSequence(query)
	results := cosineScore(e.indexer.index, terms)

	for _, result := range results {
		title, err := e.documentManager.fetchTitle(result.docId)
		if err != nil {
			return nil, err
		}
		result.title = title
	}
	return results, nil
}
