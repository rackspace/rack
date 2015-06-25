package util

import "github.com/codegangsta/cli"

// OutputFlags are global flags (i.e. flags that all commands can use) that let
// users specify the format of the output from a command.
func OutputFlags() []cli.Flag {
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

// AuthFlags are global flags (i.e. flags that all commands can use) that let
// users specify authentication parameters.
func AuthFlags() []cli.Flag {
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

// GlobalFlags returns the flags that can be used after `rack` in a command, such as
// output flags and authentication flags.
func GlobalFlags() []cli.Flag {
	gFlags := OutputFlags()
	gFlags = append(gFlags, AuthFlags()...)
	return gFlags
}
