package commands

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
)

// インデックスを作成するコマンド
var createIndexCommand = cli.Command{
	Name:   "create",
	Usage:  "create index",
	Action: createIndex,
}

func createIndex(c *cli.Context) error {

	if err := checkArgs(c, 1, exactArgs); err != nil {
		return err
	}
	dir := c.Args().Get(0)

	var files []string
	// 指定されたディレクトリ配下の.txtファイルのパスをすべて取得
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
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
		err := func() error {
			fp, err := os.Open(file)
			if err != nil {
				return err
			}
			defer fp.Close()
			title := filepath.Base(file)
			if err = engine.AddDocument(title, fp); err != nil {
				return err
			}
			log.Printf("add document to index: %s\n", title)
			return nil
		}()
		if err != nil {
			log.Println(err)
		}
	}

	return engine.Flush()
}
