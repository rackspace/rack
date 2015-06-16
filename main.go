package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/jrperritt/rack/filescommands"
	"github.com/jrperritt/rack/networkscommands"
	"github.com/jrperritt/rack/serverscommands"
	"github.com/jrperritt/rack/util"

	"github.com/codegangsta/cli"
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
	app.Usage = "An opinionated CLI for the Rackspace cloud"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:   "configure",
			Usage:  "Used to interactively create a config file for Rackspace authentication.",
			Action: configure,
		},
		{
			Name:        "servers",
			Usage:       "Used for the Servers service",
			Subcommands: serverscommands.Get(),
		},
		{
			Name:        "files",
			Usage:       "Used for the Files service",
			Subcommands: filescommands.Get(),
		},
		{
			Name:        "networks",
			Usage:       "Used for the Networks service",
			Subcommands: networkscommands.Get(),
		},
{
			Name:        "blockstorage",
			Usage:       "Used for the BlockStorage service",
			Subcommands: blockstoragecommands.Get(),
		},
	}
	app.Flags = util.GlobalFlags()
	app.BashComplete = func(c *cli.Context) {
		completeGlobals(globalOptions(app))
	}
	app.Before = func(c *cli.Context) error {
		//fmt.Printf("c.Args: %+v\n", c.Args())
		return nil
	}
	app.CommandNotFound = commandNotFound
	app.Run(os.Args)
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
	globalFlags := util.GlobalFlags()
	for _, globalFlag := range globalFlags {
		i = append(i, globalFlag)
	}

	for _, cmd := range app.Commands {
		i = append(i, cmd)
	}
	return i
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
