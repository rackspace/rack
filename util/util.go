package util

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
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

// CommonFlags are flags that all commands can use. There exists the possiblity
// of setting app-level (global) flags, but that requires a user to properly
// position them. Including these with the other command-level flags will allow
// users to include them anywhere after the last subcommand (or argument, if applicable).
func commonFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "json",
			Usage: "Return output in JSON format.",
		},
		cli.BoolFlag{
			Name:  "table",
			Usage: "Return output in tabular format. This is the default output format.",
		},
	}
}

// CommandFlags returns the flags for a given command. It takes as a parameter
// a function for returning flags specific to that command, and then appends those
// flags with flags that are valid for all commands.
func CommandFlags(f func() []cli.Flag) []cli.Flag {
	cf := commonFlags()
	return append(cf, f()...)
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

// SimpleMapsTable writes a table listing from an array of map[string]interface{}
func SimpleMapsTable(c *cli.Context, many []map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(c.App.Writer, 0, 8, 1, '\t', 0)
	// Write the header
	fmt.Fprintln(w, strings.Join(keys, "\t"))

	for _, m := range many {
		writeMapEntry(w, m, keys)
	}
	w.Flush()
}

// MetaDataPrint writes standardized metadata out
func MetaDataPrint(c *cli.Context, i interface{}, keys []string) {
	m := structs.Map(i)

	MetaDataMapPrint(c, m, keys)
}

// MetaDataMapPrint writes standardized metadata out
func MetaDataMapPrint(c *cli.Context, m map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(c.App.Writer, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, "PROPERTY\tVALUE")

	for _, key := range keys {
		val := fmt.Sprint(m[key])
		fmt.Fprintf(w, "%s\t%s\n", key, strings.Replace(val, "\n", "\n\t", -1))
	}
	w.Flush()
}

// writeMapEntry writes a table entry from a map
func writeMapEntry(w *tabwriter.Writer, m map[string]interface{}, keys []string) {
	f := []string{}
	for _, key := range keys {
		f = append(f, fmt.Sprint(m[key]))
	}
	fmt.Fprintln(w, strings.Join(f, "\t"))
}
