package util

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

// outputFlags are flags that all commands can use. There exists the possiblity
// of setting app-level (global) flags, but that requires a user to properly
// position them. Including these with the other command-level flags will allow
// users to include them anywhere after the last subcommand (or argument, if applicable).
func outputFlags() []cli.Flag {
	of := []cli.Flag{
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

	return of
}

// GlobalFlags returns the global flags for a given command.
func GlobalFlags() []cli.Flag {
	of := outputFlags()
	return of
}

// CommandFlags returns the flags for a given command. It takes as a parameter
// a function for returning flags specific to that command, and then appends those
// flags with flags that are valid for all commands.
func CommandFlags(f func() []cli.Flag, keys []string) []cli.Flag {
	//of := outputFlags(fields)
	//return append(of, f()...)
	of := f()
	if len(keys) > 0 {
		fields := make([]string, len(keys))
		for i, key := range keys {
			fields[i] = strings.Join(strings.Split(strings.ToLower(key), " "), "")
		}
		flagField := cli.StringFlag{
			Name:  "fields",
			Usage: fmt.Sprintf("Only return these comma-separated case-insensitive fields for each item in the list.\n\tChoices: %s", strings.Join(fields, ", ")),
		}
		of = append(of, flagField)
	}

	return of
}

// CompleteFlags returns the possible flags for bash completion.
func CompleteFlags(flags []cli.Flag) {
	for _, flag := range flags {
		flagName := ""
		switch flag.(type) {
		case cli.StringFlag:
			flagName = flag.(cli.StringFlag).Name
		case cli.IntFlag:
			flagName = flag.(cli.IntFlag).Name
		case cli.BoolFlag:
			flagName = flag.(cli.BoolFlag).Name
		default:
			continue
		}
		fmt.Println("--" + flagName)
	}
}

// CheckKVFlag is a function used for verifying the format of a key-value flag.
func CheckKVFlag(c *cli.Context, flagName string) map[string]string {
	kv := make(map[string]string)
	kvStrings := strings.Split(c.String(flagName), ",")
	for _, kvString := range kvStrings {
		temp := strings.Split(kvString, "=")
		if len(temp) != 2 {
			PrintError(c, ErrFlagFormatting{
				Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s for --%s.\n", kvString, flagName),
			})
		}
		kv[temp[0]] = temp[1]
	}
	return kv
}
