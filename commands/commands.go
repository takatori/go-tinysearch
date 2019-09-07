package commands

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/takatori/go-tinysearch"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	exactArgs = iota
	minArgs
	maxArgs
)

var engine *tinysearch.Engine

func Main() {

	app := cli.NewApp()
	app.Name = "tinysearch"
	app.Usage = `simple and small search engine for learning`
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		createIndexCommand,
		searchCommand,
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/tinysearch")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	engine = tinysearch.NewSearchEngine(db)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func checkArgs(context *cli.Context, expected, checkType int) error {
	var err error
	cmdName := context.Command.Name
	switch checkType {
	case exactArgs:
		if context.NArg() != expected {
			err = fmt.Errorf(
				"%s: %q requires exactly %d argument(s)",
				os.Args[0], cmdName, expected)
		}
	case minArgs:
		if context.NArg() < expected {
			err = fmt.Errorf(
				"%s: %q requires a minimum of %d argument(s)",
				os.Args[0], cmdName, expected)
		}
	case maxArgs:
		if context.NArg() > expected {
			err = fmt.Errorf(
				"%s: %q requires a maximum of %d argument(s)",
				os.Args[0], cmdName, expected)
		}
	}
	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		cli.ShowCommandHelp(context, cmdName)
		return err
	}
	return nil
}
