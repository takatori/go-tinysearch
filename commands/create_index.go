package commands

import (
	"github.com/takatori/go-tinysearch"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
)

func createIndex(c *cli.Context) error {

	if err := checkArgs(c, 1, exactArgs); err != nil {
		return err
	}
	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	engine := tinysearch.NewSearchEngine(db)

	var files []string
	dir := c.Args().Get(0)

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

			err = engine.AddDocument(file, fp)

			if err != nil {
				return err
			}
			return nil
		}()
	}

	if err := engine.Flush(); err != nil {
		return err
	}
	return nil
}
