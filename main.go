package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/jrperritt/rack/commands/blockstoragecommands"
	"github.com/jrperritt/rack/commands/filescommands"
	"github.com/jrperritt/rack/commands/networkscommands"
	"github.com/jrperritt/rack/commands/serverscommands"
	"github.com/jrperritt/rack/setup"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
)

func main() {
	cli.HelpPrinter = printHelp
	cli.CommandHelpTemplate = `NAME: {{.Name}} - {{.Usage}}{{if .Description}}

DESCRIPTION: {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
{{range .Flags}}{{flag .}}
{{end}}{{ end }}
`
	app := cli.NewApp()
	app.Name = "rack"
	app.Usage = Usage()
	app.EnableBashCompletion = true
	app.Commands = Cmds()
	app.Before = func(c *cli.Context) error {
		//fmt.Printf("c.Args: %+v\n", c.Args())
		return nil
	}
	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
}

// Usage returns, you guessed it, the usage information
func Usage() string {
	return "An opinionated CLI for the Rackspace cloud"
}

// Desc returns, you guessed it, the description
func Desc() string {
	return `Rack is an opinionated command-line tool that allows Rackspace users
to accomplish tasks in a simple, idiomatic way. It seeks to provide
flexibility through common Unix practices like piping and composability. All
commands have been tested against Rackspace's live API.`
}

// Cmds returns a list of commands supported by the tool
func Cmds() []cli.Command {
	return []cli.Command{
		{
			Name:   "init",
			Usage:  "[Linux/OS X only] Setup environment with command completion for the Bash shell.",
			Action: setup.Init,
		},
		{
			Name:   "configure",
			Usage:  "Interactively create a config file for Rackspace authentication.",
			Action: configure,
		},
		{
			Name:        "servers",
			Usage:       "Operations on cloud servers, both virtual and bare metal.",
			Subcommands: serverscommands.Get(),
		},
		{
			Name:        "files",
			Usage:       "Object storage for files and media.",
			Subcommands: filescommands.Get(),
		},
		{
			Name:        "networks",
			Usage:       "Software-defined networking.",
			Subcommands: networkscommands.Get(),
		},
		{
			Name:        "block-storage",
			Usage:       "Block-level storage, exposed as volumes to mount to host servers. Work with volumes and their associated snapshots.",
			Subcommands: blockstoragecommands.Get(),
		},
	}
}

func printHelp(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{
		"join": strings.Join,
		"flag": flag,
	}

	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func flag(flag cli.Flag) string {
	switch flag.(type) {
	case cli.StringFlag:
		flagType := flag.(cli.StringFlag)
		return fmt.Sprintf("%s\t%s", flagType.Name, flagType.Usage)
	case cli.IntFlag:
		flagType := flag.(cli.IntFlag)
		return fmt.Sprintf("%s\t%s", flagType.Name, flagType.Usage)
	case cli.BoolFlag:
		flagType := flag.(cli.BoolFlag)
		return fmt.Sprintf("%s\t%s", flagType.Name, flagType.Usage)
	case cli.StringSliceFlag:
		flagType := flag.(cli.StringSliceFlag)
		return fmt.Sprintf("%s\t%s", flagType.Name, flagType.Usage)
	}
	return ""
}
