package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/ehazlett/interlock"
)

func main() {
	app := cli.NewApp()
	app.Name = "interlock"
	app.Usage = "event driven docker plugins"
	app.Version = interlock.VERSION
	app.Email = "github.com/ehazlett/interlock"
	app.Author = "@ehazlett"

	app.Commands = appCommands

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
