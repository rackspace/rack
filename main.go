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
	app.Run(os.Args)
}
