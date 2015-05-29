package imagecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/clients"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/images"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get [flags]", util.Name, commandPrefix),
	Description: "Retreives an image",
	Action:      commandGet,
	Flags:       flagsGet(),
}

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 1)
	flavorID := c.Args()[0]
	client := clients.Create("compute")
	o, err := images.Get(client, flavorID).Extract()
	if err != nil {
		fmt.Printf("Error retreiving image [%s]: %s\n", flavorID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}

func tableGet(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetHeader([]string{"property", "value"})
	for k, v := range m {
		t.Append([]string{k, fmt.Sprint(v)})
	}
	t.Render()
}
