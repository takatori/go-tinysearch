package commands

import (
	"database/sql"
	"fmt"
	"github.com/takatori/go-tinysearch"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
)

const (
	exactArgs = iota
	minArgs
	maxArgs
)

func checkArgs(context *cli.Context, expected, checkType int) error {
	var err error
	cmdName := context.Command.Name
	switch checkType {
	case exactArgs:
		if context.NArg() != expected {
			err = fmt.Errorf("%s: %q requires exactly %d argument(s)", os.Args[0], cmdName, expected)
		}
	case minArgs:
		if context.NArg() < expected {
			err = fmt.Errorf("%s: %q requires a minimum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	case maxArgs:
		if context.NArg() > expected {
			err = fmt.Errorf("%s: %q requires a maximum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	}

	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		cli.ShowCommandHelp(context, cmdName)
		return err
	}
	return nil
}

var Commands = []cli.Command{
	{
		Name:      "create",
		Usage:     "create index",
		Description: "",
		Action: func(c *cli.Context) error {
			if err := checkArgs(c, 1, exactArgs); err != nil {
				return err
			}

			db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/tinysearch")
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()
			engine := tinysearch.NewSearchEngine(db)

			// 指定したディレクトリ配下の.txtファイルをすべて取得する
			var files []string
			root := c.Args().Get(0)

			err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
				func() error { // TODO: error handling
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

			err = engine.Flush()
			if err != nil {
				return err
			}
			return nil
		},
	},
	{
		Name:    "search",
		Aliases: []string{"s"},
		Usage:   "search for documents that match the query",
		Action: func(c *cli.Context) error {

			if err := checkArgs(c, 1, exactArgs); err != nil {
				return err
			}

			db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/tinysearch")
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()
			engine := tinysearch.NewSearchEngine(db)

			query := c.Args().Get(0)

			result, err := engine.Search(query, 10)
			if err != nil {
				return err
			}
			fmt.Println(result)

			return nil
		},
	},
}
