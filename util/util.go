package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

// Name is the name of the CLI
var Name = "rack"

// Version is the current CLI version
var Version = "0.1"

// Contains checks whether a given string is in a provided slice of strings.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// OutputFlags are flags that all commands can use. There exists the possiblity
// of setting app-level (global) flags, but that requires a user to properly
// position them. Including these with the other command-level flags will allow
// users to include them anywhere after the last subcommand (or argument, if applicable).
func OutputFlags() []cli.Flag {
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

// CheckArgNum checks that the provided number of arguments has the same
// cardinality as the expected number of arguments.
func CheckArgNum(c *cli.Context, expected int) {
	argsLen := len(c.Args())
	if argsLen != expected {
		fmt.Printf("Expected %d args but got %d\nUsage: %s\n", expected, argsLen, c.Command.Usage)
		os.Exit(1)
	}
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
