package main

import (
	"os"

	"github.com/jrperritt/rack/computecommands"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rack"
	app.Usage = "An opinionated CLI for the Rackspace cloud"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:        "compute",
			Usage:       "Used for the Compute service",
			Subcommands: computecommands.Get(),
		},
	}
	app.Run(os.Args)
}
