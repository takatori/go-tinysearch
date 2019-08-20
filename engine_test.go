package tinysearch

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// DBの初期化を行う関数
func setup() *sql.DB {

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/tinysearch")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE TABLE documents`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// インデックス構築処理のテスト
func TestCreateIndex(t *testing.T) {

	db := setup()
	defer db.Close()
	engine := NewSearchEngine(db) // ❶ 検索エンジンを初期化する

	// 指定したディレクトリ配下の.txtファイルをすべて取得する
	var files []string
	root := "testdata/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// インデックスの構築を行う
	// ファイルをひとつずつ読み込み、インデックスに追加していく
	for _, file := range files {
		func() {
			fp, err := os.Open(file)
			if err != nil {
				t.Fatalf("failed read data from %s: %v", file, err)
			}
			defer fp.Close()
			err = engine.AddDocument(file, fp) // ❷ インデックスにドキュメントを追加する
			if err != nil {
				t.Fatalf("failed to add document to index %s: %v", file, err)
			}
		}()
	}

	err = engine.Flush() // ❸ インデックスをファイルに書き出して保存
	if err != nil {
		t.Fatalf("failed to save index to file :%v", err)
	}

	// 結果の検証
	file, err := os.Open("index.json")
	if err != nil {
		t.Fatalf("failed to load index: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	expected := `{"Dictionary":{"a":[{"DocID":3,"Positions":[12],"TermFrequency":1}],"am":[{"DocID":3,"Positions":[5],"TermFrequency":1}],"as":[{"DocID":3,"Positions":[10,14],"TermFrequency":2}],"better":[{"DocID":4,"Positions":[1],"TermFrequency":1}],"do":[{"DocID":1,"Positions":[0],"TermFrequency":1},{"DocID":3,"Positions":[2],"TermFrequency":1}],"for":[{"DocID":3,"Positions":[6],"TermFrequency":1}],"good":[{"DocID":3,"Positions":[11],"TermFrequency":1}],"i":[{"DocID":3,"Positions":[4,8],"TermFrequency":2}],"if":[{"DocID":3,"Positions":[0],"TermFrequency":1}],"man":[{"DocID":3,"Positions":[13],"TermFrequency":1}],"no":[{"DocID":2,"Positions":[2],"TermFrequency":1},{"DocID":4,"Positions":[0],"TermFrequency":1}],"quarrel":[{"DocID":1,"Positions":[2],"TermFrequency":1},{"DocID":2,"Positions":[0],"TermFrequency":1}],"serve":[{"DocID":3,"Positions":[9],"TermFrequency":1}],"sir":[{"DocID":1,"Positions":[3],"TermFrequency":1},{"DocID":2,"Positions":[1,3],"TermFrequency":2},{"DocID":3,"Positions":[3],"TermFrequency":1},{"DocID":5,"Positions":[1],"TermFrequency":1}],"well":[{"DocID":5,"Positions":[0],"TermFrequency":1}],"you":[{"DocID":1,"Positions":[1],"TermFrequency":1},{"DocID":3,"Positions":[1,7,15],"TermFrequency":3}]},"DocCount":5,"DocLength":{"1":4,"2":4,"3":16,"4":2,"5":2}}`

	if string(bytes) != expected { // TODO: ちゃんと比較する
		t.Fatalf("failed to create index")
	}

	// [For Search]
	// when
	query := "Quarrel, sir."
	actual, err := engine.Search(query, 5)
	if err != nil {
		t.Fatalf("failed search: %v", err)
	}

	// then
	expectedSearchResult := []*SearchResult{
		{2, 0.66, "testdata/romeo_and_juliet_2.txt"},
		{1, 0.59, "testdata/romeo_and_juliet_1.txt"},
		{5, 0.21, "testdata/romeo_and_juliet_5.txt"},
		{3, 0.03, "testdata/romeo_and_juliet_3.txt"},
	}

	for i := range expectedSearchResult {
		if actual[i].docID != expectedSearchResult[i].docID {
			t.Fatalf("\ngot:\n%v\nwant:\n%v\n", actual, expected)
		}
	}
}
