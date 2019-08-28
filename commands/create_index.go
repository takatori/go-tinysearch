package commands

import (
	"database/sql"
	"github.com/takatori/go-tinysearch"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
)

func createIndex(c *cli.Context) error {

	if err := checkArgs(c, 1, exactArgs); err != nil {
		return err
	}
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/tinysearch")
	if err != nil {
		return err
	}
	defer db.Close()

	engine := tinysearch.NewSearchEngine(db)

	var files []string
	dir := c.Args().Get(0)

	// 指定されたディレクトリ配下の.txtファイルのパスをすべて取得
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".txt" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		func() error {
			fp, err := os.Open(file)
			if err != nil {
				return err
			}
			defer fp.Close()
			// ドキュメントの追加
			if err = engine.AddDocument(file, fp); err != nil {
				return err
			}
			log.Printf("add document to index: %s\n", file)
			return nil
		}()
	}

	// インデックスを永続化
	if err := engine.Flush(); err != nil {
		return err
	}
	return nil
}
