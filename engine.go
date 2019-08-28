package tinysearch

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// 検索エンジン
type Engine struct {
	tokenizer     *Tokenizer     // トークンを分割する
	indexer       *Indexer       // インデックスを作成する
	documentStore *DocumentStore // ドキュメントを管理する
	indexFile     string         // インデックスファイルのパス
}

// 検索エンジンを作成する処理
func NewSearchEngine(db *sql.DB) *Engine {

	tokenizer := NewTokenizer()
	indexer := NewIndexer(tokenizer)
	documentStore := NewDocumentStore(db)
	indexFilePath, ok := os.LookupEnv("INDEX_FILE_PATH")
	if !ok {
		indexFilePath = "index.json"
	}

	return &Engine{
		tokenizer:     tokenizer,
		indexer:       indexer,
		documentStore: documentStore,
		indexFile:     indexFilePath,
	}
}

// インデックスにドキュメントを追加する
func (e *Engine) AddDocument(title string, reader io.Reader) error {
	id, err := e.documentStore.save(title) // ❶ タイトルを保存しドキュメントIDを発行する
	if err != nil {
		return err
	}
	e.indexer.update(id, reader) // ❷ インデックスを更新する
	return nil
}

// インデックスをファイルに書き出す
func (e *Engine) Flush() error {

	bytes, err := json.Marshal(e.indexer.index)
	if err != nil {
		return err
	}

	file, err := os.Create(e.indexFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}
	return writer.Flush()

	// TODO: indexer.indexを空にする？
}

// 検索を実行する
func (e *Engine) Search(query string, k int) ([]*SearchResult, error) {

	// クエリをトークンに分割
	terms := e.tokenizer.TextToWordSequence(query)

	// インデックスを読み込む
	file, err := os.Open(e.indexFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	idx := NewIndex()
	if err := json.Unmarshal(bytes, idx); err != nil {
		return nil, err
	}

	// 検索を実行
	searcher := NewSearcher(idx)
	docs := searcher.searchTopK(terms, k)

	// タイトルを取得
	results := make([]*SearchResult, 0, k)
	for _, result := range docs.scoreDocs {
		title, err := e.documentStore.fetchTitle(result.docID)
		if err != nil {
			return nil, err
		}
		results = append(results, &SearchResult{
			result.docID, result.score, title,
		})
	}
	return results, nil
}

// 検索結果を格納する構造体
type SearchResult struct {
	DocID DocumentID
	Score float64
	Title string
}

// String print searchTopK result info
func (r *SearchResult) String() string {
	return fmt.Sprintf("{DocID: %v, Score: %v, Title: %v}",
		r.DocID, r.Score, r.Title)
}
