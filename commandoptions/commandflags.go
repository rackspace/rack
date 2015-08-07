package commandoptions

import (
	"fmt"
	"strings"

	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
)

// CommandFlags returns the flags for a given command. It takes as a parameter
// a function for returning flags specific to that command, and then appends those
// flags with flags that are valid for all commands.
func CommandFlags(f func() []cli.Flag, keys []string) []cli.Flag {
	of := f()
	if len(keys) > 0 {
		fields := make([]string, len(keys))
		for i, key := range keys {
			fields[i] = strings.Join(strings.Split(strings.ToLower(key), " "), "-")
		}
		flagFields := cli.StringFlag{
			Name:  "fields",
			Usage: fmt.Sprintf("[optional] Only return these comma-separated case-insensitive fields.\nChoices: %s", strings.Join(fields, ", ")),
		}
		of = append(of, flagFields)
	}
	of = append(of, GlobalFlags()...)

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
		case cli.StringSliceFlag:
			flagName = flag.(cli.StringSliceFlag).Name
		default:
			continue
		}
		fmt.Println("--" + flagName)
	}
}
