package util

import "github.com/jrperritt/rack/internal/github.com/codegangsta/cli"

// GlobalFlags returns the flags that can be used after `rack` in a command, such as
// output flags and authentication flags.
func GlobalFlags() []cli.Flag {
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
		cli.StringFlag{
			Name:  "output",
			Usage: "Format in which to return output. Options: json, csv, table. Default is 'table'.",
		},
		cli.BoolFlag{
			Name:  "no-cache",
			Usage: "Don't get or set authentication credentials in the rack cache.",
		},
		cli.StringFlag{
			Name:  "log",
			Usage: "Print debug information from the command. Options are: debug, info",
		},
	}
}
