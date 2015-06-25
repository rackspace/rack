package main

import (
	"fmt"
	"os"

	"github.com/jrperritt/rack/filescommands"
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
			Name:        "files",
			Usage:       "Used for the Files service",
			Subcommands: filescommands.Get(),
		},
	}
	app.Flags = util.GlobalFlags()
	app.BashComplete = func(c *cli.Context) {
		completeGlobals(globalOptions(app))
	}
	app.Run(os.Args)
}

// completeGlobals returns the options for completing global flags.
func completeGlobals(vals []interface{}) {
	for _, val := range vals {
		switch val.(type) {
		case cli.StringFlag:
			fmt.Println("--" + val.(cli.StringFlag).Name)
		case cli.IntFlag:
			fmt.Println("--" + val.(cli.IntFlag).Name)
		case cli.BoolFlag:
			fmt.Println("--" + val.(cli.BoolFlag).Name)
		case cli.Command:
			fmt.Println(val.(cli.Command).Name)
		default:
			continue
		}
	}
}

// globalOptions returns the options (flags and commands) that can be used after
// `rack` in a command. For example, `rack --json`, `rack servers`, and
// `rack --json servers` are all legitimate command prefixes.
func globalOptions(app *cli.App) []interface{} {
	var i []interface{}
	globalFlags := util.GlobalFlags()
	for _, globalFlag := range globalFlags {
		i = append(i, globalFlag)
	}

	for _, cmd := range app.Commands {
		i = append(i, cmd)
	}
	return i
}
