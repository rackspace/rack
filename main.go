package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rackspace/rack/commands/blockstoragecommands"
	"github.com/rackspace/rack/commands/filescommands"
	"github.com/rackspace/rack/commands/networkscommands"
	"github.com/rackspace/rack/commands/orchestrationcommands"
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
	return "Command-line interface to manage Rackspace Cloud resources"
}

// Desc returns, you guessed it, the description
func Desc() string {
	return `The rack CLI manages authentication, configures a local setup, and provides workflows for operations on Rackspace Cloud resources`
}

// Cmds returns a list of commands supported by the tool
func Cmds(app cli.App) []cli.Command {
	//isAdmin := util.IsAdmin()

	return []cli.Command{
		{
			Name:   "configure",
			Usage:  "Interactively create a config file for Rackspace authentication",
			Action: configure,
		},
		{
			Name: "init",
			Usage: "Enable tab for command completion.\n" +
				"\tFor Linux and OS X, creates the `rack` man page and sets up\n" +
				"\tcommand completion for the Bash shell. Run `man ./rack.1` to\n" +
				"\tview the generated man page.\n" +
				"\tFor Windows, creates a `posh_autocomplete.ps1` file in the\n" +
				"\t`$HOME/.rack` directory. You must run the file to set up\n" +
				"\tcommand completion\n",
			Action: func(c *cli.Context) {
				setup.Init(c)
				man()
			},
		},
		{
			Name:  "version",
			Usage: "Print the version of this binary",
			Action: func(c *cli.Context) {
				fmt.Fprintf(c.App.Writer, "%v version %v\ncommit: %v\n", c.App.Name, util.Version, util.Commit)
			},
		},
		{
			Name:        "profile",
			Usage:       "Used to perform operations on user profiles",
			Subcommands: profileCommandsGet(isAdmin),
		},
		{
			Name:        "servers",
			Usage:       "Operations on cloud servers, both virtual and bare metal",
			Subcommands: serverscommands.Get(),
		},
		{
			Name:        "files",
			Usage:       "Object storage for files and media",
			Subcommands: filescommands.Get(),
		},
		{
			Name:        "networks",
			Usage:       "Software-defined networking",
			Subcommands: networkscommands.Get(),
		},
		{
			Name: "block-storage",
			Usage: strings.Join([]string{"Block-level storage, exposed as volumes to mount to",
				"\thost servers. Work with volumes and their associated snapshots"}, "\n"),
			Subcommands: blockstoragecommands.Get(),
		},
		{
			Name:        "orchestration",
			Usage:       "Use a template language to orchestrate Rackspace cloud services",
			Subcommands: orchestrationcommands.Get(),
		},
	}
}
