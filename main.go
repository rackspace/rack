package main

import (
	"os"

	"github.com/jrperritt/rack/serverscommands"
	"github.com/jrperritt/rack/util"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rack"
	app.Usage = "An opinionated CLI for the Rackspace cloud"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:        "servers",
			Usage:       "Used for the Servers service",
			Subcommands: serverscommands.Get(),
		},
	}
	app.Flags = util.GlobalFlags()
	app.BashComplete = func(c *cli.Context) {
		util.CompleteGlobals(util.GlobalOptions(app))
	}
	app.Run(os.Args)
}
