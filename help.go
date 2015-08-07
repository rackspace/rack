package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
)

var appHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}} <command> <subcommand> <action> [FLAGS]
   {{if .Version}}
VERSION:
   {{.Version}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{wrap .Usage}}
   {{end}}{{end}}
`

var commandHelpTemplate = `NAME: {{.Name}} - {{.Usage}}{{if .Description}}

DESCRIPTION: {{.Description}}{{end}}{{if .Flags}}

COMMAND FLAGS:
{{range .Flags}}{{if isNotGlobalFlag .}}{{flag .}}
{{end}}{{end}}

GLOBAL FLAGS:
{{range .Flags}}{{if isGlobalFlag .}}{{flag .}}
{{end}}{{end}}{{ end }}
`

var subcommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}}{{if eq (len (split .Name " ")) 2}} <subcommand>{{end}} <action> [FLAGS]
{{if eq (len (split .Name " ")) 2}}SUBCOMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{wrap .Usage}}
   {{end}}{{else}}ACTIONS:
	 {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
	 {{end}}{{end}}
`

func printHelp(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{
		"split":           strings.Split,
		"join":            strings.Join,
		"isGlobalFlag":    isGlobalFlag,
		"isNotGlobalFlag": isNotGlobalFlag,
		"flag":            flag,
		"wrap":            wrap,
	}

	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func isGlobalFlag(cliflag cli.Flag) bool {
	globalFlags := commandoptions.GlobalFlags()
	for _, globalFlag := range globalFlags {
		if globalFlag == cliflag {
			return true
		}
	}
	return false
}

func isNotGlobalFlag(cliflag cli.Flag) bool {
	globalFlags := commandoptions.GlobalFlags()
	for _, globalFlag := range globalFlags {
		if globalFlag == cliflag {
			return false
		}
	}
	return true
}

func flag(cliflag cli.Flag) string {
	var flagString string
	switch cliflag.(type) {
	case cli.StringFlag:
		flagType := cliflag.(cli.StringFlag)
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), wrap(flagType.Usage))
	case cli.IntFlag:
		flagType := cliflag.(cli.IntFlag)
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), wrap(flagType.Usage))
	case cli.BoolFlag:
		flagType := cliflag.(cli.BoolFlag)
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), wrap(flagType.Usage))
	case cli.StringSliceFlag:
		flagType := cliflag.(cli.StringSliceFlag)
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), wrap(flagType.Usage))
	}
	return flagString
}

func wrap(text string) string {
	textSlice := strings.Split(text, "\n")
	var wrapSlice []string
	for _, line := range textSlice {
		for j := 0; j < len(line); j += 55 {
			//fmt.Printf("j:%d, j+55:%d, len(text):%d\n", j, j+55, len(text))
			offset := 0
			end := j + 55
			if end > len(line) {
				end = len(line)
			} else {
				for {
					//fmt.Printf("text[end:end+offset]: %s\n", string(text[end:end+offset]))
					if end+offset == len(line) || line[end+offset] == ' ' {
						break
					}
					offset++
				}
				end += offset
			}
			//fmt.Printf("line: %s\n", line[j:end])
			wrapSlice = append(wrapSlice, line[j:end])
			j += offset + 1
		}
	}
	return strings.Join(wrapSlice, "\n\t")
}
