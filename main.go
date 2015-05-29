package main

import (
	"os"

	"github.com/jrperritt/rackcli/computecommands"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rackcli"
	app.Usage = "An opinionated CLI for the Rackspace cloud"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:        "compute",
			Usage:       "Used for the Compute service",
			Subcommands: computecommands.Get(),
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "The format for the output. Options are json and table. Default is table.",
		},
	}
	app.Run(os.Args)
}
