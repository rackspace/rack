package main

import (
	"fmt"
	"html/template"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
)

var commandHelpTemplate = `NAME: {{.Name}} - {{.Usage}}{{if .Description}}

DESCRIPTION: {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
{{range .Flags}}{{flag .}}
{{end}}{{ end }}
`

var appHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}} <command> <subcommand> <action> [OPTIONS]
   {{if .Version}}
VERSION:
   {{.Version}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}
`

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
