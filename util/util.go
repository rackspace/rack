package util

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

// Name is the name of the CLI
var Name = "rackcli"

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
func CompleteFlags(f func() []cli.Flag) {
	flags := f()
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
