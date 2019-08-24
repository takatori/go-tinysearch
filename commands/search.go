package commands

import (
	"fmt"
	"github.com/takatori/go-tinysearch"
	"github.com/urfave/cli"
	"log"
)

func search(c *cli.Context) error {

	if err := checkArgs(c, 1, exactArgs); err != nil {
		return err
	}

	db, err := db()
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
}
