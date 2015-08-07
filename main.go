package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rackspace/rack/commands/blockstoragecommands"
	"github.com/rackspace/rack/commands/filescommands"
	"github.com/rackspace/rack/commands/networkscommands"
	"github.com/rackspace/rack/commands/serverscommands"
	"github.com/rackspace/rack/setup"
	"github.com/rackspace/rack/util"

	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
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
	return "An opinionated CLI for the Rackspace cloud"
}

// Desc returns, you guessed it, the description
func Desc() string {
	return `Rack is an opinionated command-line tool that allows Rackspace users
to accomplish tasks in a simple, idiomatic way. It seeks to provide
flexibility through common Unix practices like piping and composability. All
commands have been tested against Rackspace's live API.`
}

// Cmds returns a list of commands supported by the tool
func Cmds() []cli.Command {
	return []cli.Command{
		{
			Name: "init",
			Usage: strings.Join([]string{"For Linux and OS X, creates the `rack` man page and sets up command completion for the Bash shell.",
				"For Windows, creates a `posh_autocomplete.ps1` file in the `$HOME/.rack` directory. Windows users must run the file to set up command completion."}, "\n"),
			Action: func(c *cli.Context) {
				setup.Init(c)
				man()
			},
		},
		{
			Name:   "configure",
			Usage:  "Interactively create a config file for Rackspace authentication.",
			Action: configure,
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
