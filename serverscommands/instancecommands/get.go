package instancecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/olekukonko/tablewriter"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get %s [optional flags]", util.Name, commandPrefix, idOrNameUsage),
	Description: "Retrieves an existing server",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet))
	},
}

func flagsGet() []cli.Flag {
	return idAndNameFlags
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	serverID := idOrName(c, client)
	o, err := servers.Get(client, serverID).Extract()
	if err != nil {
		fmt.Printf("Error retrieving server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}

func tableGet(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{"property", "value"})
	keys := []string{"ID", "Name", "Status", "Created", "Updated", "Image", "Flavor", "Public IPv4", "Public IPv6", "Private IPv4", "KeyName"}

	mungeServerMap(m)

	for _, key := range keys {
		t.Append([]string{key, fmt.Sprint(m[key])})
	}
	t.Render()
}
