package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/takatori/go-tinysearch"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
)

// flag parse
// インデックスを作るコマンド
// 検索するコマンド
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

func main() {

	app := cli.NewApp()
	app.Name = "tinysearch"
	app.Usage = `simple and small search engine for learning`

	app.Commands = []cli.Command{
		{
			Name:      "create",
			Aliases:   []string{"c"},
			Usage:     "create index",
			ArgsUsage: ``,
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
				fmt.Println(c.Args())
				fmt.Println(root)
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

				// インデックスの構築を行う
				// ファイルをひとつずつ読み込み、インデックスに追加していく
				for _, file := range files {
					func() error { // TODO: error handling
						fp, err := os.Open(file)
						if err != nil {
							return err
						}
						defer fp.Close()
						err = engine.AddDocument(file, fp) // ❷ インデックスにドキュメントを追加する
						if err != nil {
							return err
						}
						return nil
					}()
				}

				err = engine.Flush() // ❸ インデックスをファイルに書き出して保存
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
