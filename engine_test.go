package tinysearch

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
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

	title := "test doc"
	body := "Do you quarrel, sir?"

	// ❷ インデックスにドキュメントを追加する
	err := engine.AddDocument("test doc", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to add document to index %s: %v", title, err)
	}
	// ❸ インデックスをファイルに書き出して永続化
	err = engine.Flush()
	if err != nil {
		t.Fatalf("failed to save index to file :%v", err)
	}

	// 以下は、検証用のコード
	file, err := os.Open("index.json")
	if err != nil {
		t.Fatalf("failed to load index: %v", err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	got := NewIndex()
	if err := json.Unmarshal(bytes, got); err != nil {
		t.Fatalf("failed to unmarshal idxStr: %v", err)
	}

	wantStr := `
{
    "TotalDocsCount":1,
    "Dictionary":{
	"do":     [{"DocID":1,"Positions":[0],"TermFrequency":1}],
	"quarrel":[{"DocID":1,"Positions":[2],"TermFrequency":1}],
	"sir":    [{"DocID":1,"Positions":[3],"TermFrequency":1}],
	"you":    [{"DocID":1,"Positions":[1],"TermFrequency":1}]}
}
`
	want := NewIndex()
	if err := json.Unmarshal([]byte(wantStr), want); err != nil {
		t.Fatalf("failed to unmarshal idxStr: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("failed to create index\n got : %v\nwant: %v\n",
			got, want)
	}
}

func TestSearch(t *testing.T) {

}

/*
	// [For Search]
	// when
	query := "Quarrel, sir."
	actual, err := engine.Search(query, 5) // TODO: 検索に使用するインデックスファイルを指定できるようにする?
	if err != nil {
		t.Fatalf("failed searchTopK: %v", err)
	}

	// then
	expectedSearchResult := []*SearchResult{
		{2, 0.66, "testdata/romeo_and_juliet_2.txt"},
		{1, 0.59, "testdata/romeo_and_juliet_1.txt"},
		{5, 0.21, "testdata/romeo_and_juliet_5.txt"},
		{3, 0.03, "testdata/romeo_and_juliet_3.txt"},
	}

	for i := range expectedSearchResult {
		if actual[i].DocumentID != expectedSearchResult[i].DocumentID {
			t.Fatalf("\ngot:\n%v\nwant:\n%v\n", actual, expected)
		}
	}*/
