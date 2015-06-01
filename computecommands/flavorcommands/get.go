package flavorcommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/flavors"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get <flavorID> [flags]", util.Name, commandPrefix),
	Description: "Retreives a flavor",
	Action:      commandGet,
	Flags:       flagsGet(),
}

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 1)
	flavorID := c.Args()[0]
	client := auth.NewClient("compute")
	o, err := flavors.Get(client, flavorID).Extract()
	if err != nil {
		fmt.Printf("Error retreiving flavor [%s]: %s\n", flavorID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}

func tableGet(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{"property", "value"})
	keys := []string{"ID", "Name", "Disk", "RAM", "RxTxFactor", "Swap", "VCPUs"}
	for _, key := range keys {
		t.Append([]string{key, fmt.Sprint(m[key])})
	}
	t.Render()
}
