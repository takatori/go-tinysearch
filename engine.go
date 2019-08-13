package tinysearch

import (
	"database/sql"
)

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

// 指定したパスのファイルから、インデックスを更新する
func (e *Engine) AddDocument(title, document string) error {

	id, err := e.documentManager.saveDocument(title)
	if err != nil {
		return err
	}

	// todo: stringを持ち回さない
	// todo: errorハンドリング
	return e.indexManager.updatePostingsList(id, document)
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
