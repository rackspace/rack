package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// outputFlags are global flags (i.e. flags that all commands can use) that let
// users specify the format of the output from a command.
func outputFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "json",
			Usage: "Return output in JSON format.",
		},
		cli.BoolFlag{
			Name:  "table",
			Usage: "Return output in tabular format. This is the default output format.",
		},
		cli.BoolFlag{
			Name:  "csv",
			Usage: "Return output in csv format.",
		},
	}
}

// globalFlags returns the flags that can be used after `rack` in a command, such as
// `--json`, `--csv`, and `--table`.
func globalFlags() []cli.Flag {
	outputFlags := outputFlags()
	return outputFlags
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
	globalFlags := globalFlags()
	for _, globalFlag := range globalFlags {
		i = append(i, globalFlag)
	}

	for _, cmd := range app.Commands {
		i = append(i, cmd)
	}
	return i
}
