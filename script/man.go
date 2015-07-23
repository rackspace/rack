package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	r "github.com/jrperritt/rack"
	"github.com/jrperritt/rack/commandoptions"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/util"
)

func main() {

	content := fmt.Sprintln(`.\" Manpage for rack`)
	content += fmt.Sprintln(`.\" Contact sdk-support@rackspace.com to correct errors or typos`)
	content += fmt.Sprintf(`.TH man 1 "%s" "%s" "rack man page"`+"\n", time.Now().Format("06 May 2010"), util.Version)
	content += fmt.Sprintln(`.SH NAME`)
	content += fmt.Sprintf(`rack \- %s`+"\n", r.Usage())
	content += fmt.Sprintln(`.SH SYNOPSIS`)
	content += fmt.Sprintln("rack [GLOBALS] command subcommand [OPTIONS]")
	content += fmt.Sprintln(`.SH DESCRIPTION`)
	content += fmt.Sprintf("%s\n\n\n", r.Desc())

	content += fmt.Sprintln("The following global options are available:")
	for _, flag := range commandoptions.GlobalFlags() {
		content += fmt.Sprintln(".TP")
		name, usage := parseFlag(flag)
		if name != "" && usage != "" {
			content += fmt.Sprintf(`\fB\-\-%s\fR`+"\n", strings.Replace(name, "-", `\-`, -1))
			content += fmt.Sprintln(usage)
		} else {
			content += fmt.Sprintln(flag.String())
		}
	}

	content += fmt.Sprintln(`.SH TOP-LEVEL COMMANDS`)
	for _, cmd := range r.Cmds() {
		if len(cmd.Subcommands) > 0 {
			continue
		}
		content += fmt.Sprintln(".TP")
		content += fmt.Sprintf(`\fB%s\fR`+"\n", cmd.Name)
		content += fmt.Sprintln(cmd.Usage)
	}

	for _, cmd := range r.Cmds() {
		if len(cmd.Subcommands) == 0 {
			continue
		}
		content += fmt.Sprintf(`.SH %s COMMANDS`+"\n", strings.ToUpper(cmd.Name))
		content += fmt.Sprintln(cmd.Usage)
		for _, serviceCmd := range cmd.Subcommands {
			name := strings.ToUpper(serviceCmd.Name)
			content += fmt.Sprintf(`.SS "\s-1%s COMMANDS\s0"`+"\n", name)
			content += fmt.Sprintf(`.IX Subsection "%s"`+"\n", name)
			for _, resourceCmd := range serviceCmd.Subcommands {
				content += fmt.Sprintf(`.IP "\fB%s\fR"`+"\n", resourceCmd.Usage)
				content += fmt.Sprintf(`.IX Item "%s"`+"\n", resourceCmd.Usage)
				content += fmt.Sprintln(resourceCmd.Description)
			}
		}
	}

	content += fmt.Sprintln(".SH BUGS")
	content += fmt.Sprintln("See https://github.com/jrperritt/rack/issues")

	ioutil.WriteFile("rack.1", []byte(content), 0755)
}

func parseFlag(f cli.Flag) (string, string) {
	name, usage := "", ""

	if f, ok := f.(cli.StringFlag); ok {
		name = f.Name
		usage = f.Usage
	}
	if f, ok := f.(cli.BoolFlag); ok {
		name = f.Name
		usage = f.Usage
	}
	if f, ok := f.(cli.IntFlag); ok {
		name = f.Name
		usage = f.Usage
	}

	return name, usage
}

func flagUsage(f cli.Flag) string {
	return ""
}
