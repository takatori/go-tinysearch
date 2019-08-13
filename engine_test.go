package tinysearch

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// DBの初期化を行う
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

func TestEngine(t *testing.T) {

	// given
	db := setup()
	defer db.Close()
	engine := NewSearchEngine(db)

	// 指定したパスのディレクトリをすべて読む
	var files []string
	root := "testdata/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".txt" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// インデックス構築
	for _, file := range files {
		func() {
			fp, err := os.Open(file)
			if err != nil {
				t.Fatalf("failed read data from %s: %v", file, err)
			}
			defer fp.Close()
			if err = engine.AddDocument(file, fp); err != nil {
				t.Fatalf("failed to add document to index %s: %v", file, err)
			}
		}()
	}

	// when
	query := "Quarrel, sir."
	actual, err := engine.Search(query, 5)
	if err != nil {
		t.Fatalf("failed search: %v", err)
	}

	// then
	expected := []*SearchResult{
		{2, 0.66, "testdata/romeo_and_juliet_2.txt"},
		{1, 0.59, "testdata/romeo_and_juliet_1.txt"},
		{5, 0.21, "testdata/romeo_and_juliet_5.txt"},
		{3, 0.03, "testdata/romeo_and_juliet_3.txt"},
	}

	for i := range expected {
		if actual[i].docId != expected[i].docId {
			t.Fatalf("\ngot:\n%v\nwant:\n%v\n", actual, expected)
		}
	}
}
