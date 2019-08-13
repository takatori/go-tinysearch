package tinysearch

import (
	"database/sql"
	"io"
)

// 検索エンジン
type Engine struct {
	indexManager    *IndexManager
	documentManager *DocumentManager
}

// NewSearchEngine(db) create a search engine.
// 検索エンジンを作成する
func NewSearchEngine(db *sql.DB) *Engine {
	return &Engine{
		NewIndexManager(),
		NewDocumentManager(db),
	}
}

// インデックスにドキュメントを追加する
func (e *Engine) AddDocument(title string, reader io.Reader) error {

	id, err := e.documentManager.save(title)
	if err != nil {
		return err
	}

	e.indexManager.update(id, reader)
	return nil
}

// 検索を実行する
func (e *Engine) Search(query string) []*SearchResult {
	terms := TextToWordSequence(query)
	results := cosineScore(e.indexManager.index, terms)

	for _, result := range results {
		title, _ := e.documentManager.fetchTitle(result.docId) // TODO: errorハンドリング
		result.title = title
	}
	return results
}
