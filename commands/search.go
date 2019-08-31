package commands

import (
	"database/sql"
	"fmt"
	"github.com/takatori/go-tinysearch"
	"github.com/urfave/cli"
	"log"
	"strings"
)

func search(c *cli.Context) error {

	if err := checkArgs(c, 1, exactArgs); err != nil {
		return err
	}
	query := c.Args().Get(0)

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/tinysearch")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	engine := tinysearch.NewSearchEngine(db)

	result, err := engine.Search(query, 100) // TODO: option
	if err != nil {
		return err
	}

	printResult(result)
	return nil
}

func printResult(results []*tinysearch.SearchResult) {

	if len(results) == 0 {
		fmt.Println("0 match!!")
		return
	}

	strs := make([]string, len(results))
	for i, result := range results {
		strs[i] = fmt.Sprintf("rank: %3d, score: %4f, title: %s",
			i+1, result.Score, result.Title)
	}
	fmt.Println(strings.Join(strs, "\n"))
}
