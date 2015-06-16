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

// authFlags are global flags (i.e. flags that all commands can use) that let
// users specify authentication parameters.
func authFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "username",
			Usage: "The username with which to authenticate.",
		},
		cli.StringFlag{
			Name:  "apikey",
			Usage: "The API key with which to authenticate.",
		},
		cli.StringFlag{
			Name:  "authurl",
			Usage: "The endpoint to which authenticate.",
		},
		cli.StringFlag{
			Name:  "region",
			Usage: "The region to which authenticate.",
		},
		cli.StringFlag{
			Name:  "profile",
			Usage: "The config file profile to use for authentication.",
		},
		cli.BoolFlag{
			Name:  "no-cache",
			Usage: "Don't get or set authentication credentials in the rack cache.",
		},
	}
}

// globalFlags returns the flags that can be used after `rack` in a command, such as
// output flags and authentication flags.
func globalFlags() []cli.Flag {
	gFlags := outputFlags()
	gFlags = append(gFlags, authFlags()...)
	return gFlags
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
