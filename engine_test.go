package tinysearch

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

var testDB *sql.DB

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
	// TODO: _index_dataディレクトリ削除
	return db
}

func TestMain(m *testing.M) {
	testDB = setup()
	defer testDB.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}

// インデックス構築処理のテスト
func TestCreateIndex(t *testing.T) {

	engine := NewSearchEngine(testDB) // ❶ 検索エンジンを初期化する

	title := "test doc"
	body := strings.NewReader("Do you quarrel, sir?")

	// ❷ インデックスにドキュメントを追加する
	if err := engine.AddDocument("test doc", body); err != nil {
		t.Fatalf("failed to add document to index %s: %v", title, err)
	}
	// ❸ インデックスをファイルに書き出して永続化
	if err := engine.Flush(); err != nil {
		t.Fatalf("failed to save index to file :%v", err)
	}
	type testCase struct {
		file        string
		postingsStr string
	}

	// TODO: 増やす
	testCases := []testCase{
		{"_index_data/_0.dc", "1"},
		{"_index_data/do.json", `[{"DocID":1,"Positions":[0],"TermFrequency":1}]`},
		{"_index_data/quarrel.json", `[{"DocID":1,"Positions":[2],"TermFrequency":1}]`},
		{"_index_data/sir.json", `[{"DocID":1,"Positions":[3],"TermFrequency":1}]`},
		{"_index_data/you.json", `[{"DocID":1,"Positions":[1],"TermFrequency":1}]`},
	}

	for _, testCase := range testCases {
		func() {
			file, err := os.Open(testCase.file)
			if err != nil {
				t.Fatalf("failed to load index: %v", err)
			}
			defer file.Close()
			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatalf("failed to load index: %v", err)
			}
			actual := string(bytes)
			expected := testCase.postingsStr
			if actual != expected {
				t.Errorf("failed to create index\n got : %v\nwant: %v\n", actual, expected)
			}
		}()
	}
}

func TestSearch(t *testing.T) {

	engine := NewSearchEngine(testDB)

	query := "Quarrel, sir."
	actual, err := engine.Search(query, 5)
	if err != nil {
		t.Fatalf("failed searchTopK: %v", err)
	}

	// then
	expected := []*SearchResult{
		{1, 0, "test doc"},
	}

	for !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\ngot:\n%v\nwant:\n%v\n", actual, expected)
	}

}
