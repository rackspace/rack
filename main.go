package main

import (
	"fmt"
	"os"

	"github.com/jrperritt/rack/commands/blockstoragecommands"
	"github.com/jrperritt/rack/commands/filescommands"
	"github.com/jrperritt/rack/commands/networkscommands"
	"github.com/jrperritt/rack/commands/serverscommands"
	"github.com/jrperritt/rack/setup"
	"github.com/jrperritt/rack/util"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
)

func main() {
	cli.HelpPrinter = printHelp
	cli.AppHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate
	cli.SubcommandHelpTemplate = subcommandHelpTemplate
	app := cli.NewApp()
	app.Name = "rack"
	app.Version = fmt.Sprintf("%v version %v\n   commit: %v\n", app.Name, util.Version, util.Commit)
	app.Usage = Usage()
	app.HideVersion = true
	app.EnableBashCompletion = true
	app.Commands = Cmds()
	app.Before = func(c *cli.Context) error {
		//fmt.Printf("c.Args: %+v\n", c.Args())
		return nil
	}
	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
}

// Usage returns, you guessed it, the usage information
func Usage() string {
	return "Command-line interface to manage Rackspace Cloud resources"
}

// Desc returns, you guessed it, the description
func Desc() string {
	return `The rack CLI manages authentication, configures a local setup, and provides workflows for operations on Rackspace cloud resources.`
}

// Cmds returns a list of commands supported by the tool
func Cmds() []cli.Command {
	return []cli.Command{
		{
			Name:   "configure",
			Usage:  "Interactively create a config file for Rackspace authentication.",
			Action: configure,
		},
		{
			Name:  "init",
			Usage: "[Linux/OS X only] Creates the rack man page and sets up command completion for the Bash shell.",
			Action: func(c *cli.Context) {
				setup.Init(c)
				man()
			},
		},
		{
			Name:  "version",
			Usage: "Print the version of this binary.",
			Action: func(c *cli.Context) {
				fmt.Fprintf(c.App.Writer, "%v version %v\ncommit: %v\n", c.App.Name, util.Version, util.Commit)
			},
		},
		{
			Name:        "servers",
			Usage:       "Operations on cloud servers, both virtual and bare metal.",
			Subcommands: serverscommands.Get(),
		},
		{
			Name:        "files",
			Usage:       "Object storage for files and media.",
			Subcommands: filescommands.Get(),
		},
		{
			Name:        "networks",
			Usage:       "Software-defined networking.",
			Subcommands: networkscommands.Get(),
		},
		{
			Name:        "block-storage",
			Usage:       "Block-level storage, exposed as volumes to mount to host servers. Work with volumes and their associated snapshots.",
			Subcommands: blockstoragecommands.Get(),
		},
	}
}
