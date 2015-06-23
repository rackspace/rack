package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/blockstoragecommands"
	"github.com/jrperritt/rack/serverscommands"
)

func main() {
	app := cli.NewApp()
	app.Name = "rack"
	app.Usage = "An opinionated CLI for the Rackspace cloud"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:   "configure",
			Usage:  "Used to interactively create a config file for Rackspace authentication.",
			Action: configure,
		},
		{
			Name:        "servers",
			Usage:       "Used for the Servers service",
			Subcommands: serverscommands.Get(),
		},
		{
			Name:        "blockstorage",
			Usage:       "Used for the BlockStorage service",
			Subcommands: blockstoragecommands.Get(),
		},
	}
	app.Flags = globalFlags()
	app.BashComplete = func(c *cli.Context) {
		completeGlobals(globalOptions(app))
	}
	app.Run(os.Args)
}
