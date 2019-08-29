package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/takatori/go-tinysearch/commands"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "tinysearch"
	app.Usage = `simple and small search engine for learning`
	app.Version = "0.0.1"
	app.Commands = commands.Commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
