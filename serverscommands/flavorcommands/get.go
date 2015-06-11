package flavorcommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/flavors"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get <flavorID> [flags]", util.Name, commandPrefix),
	Description: "Retreives a flavor",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet))
	},
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
	keys := []string{"ID", "Name", "Disk", "RAM", "RxTxFactor", "Swap", "VCPUs"}
	f := func() map[string]interface{} {
		return structs.Map(i)
	}
	output.MetadataTable(c, &f, keys)
}
