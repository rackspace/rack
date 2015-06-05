package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var update = cli.Command{
	Name:        "update",
	Usage:       fmt.Sprintf("%s %s update %s [optional flags]", util.Name, commandPrefix, idOrNameUsage),
	Description: "Updates an existing server",
	Action:      commandUpdate,
	Flags:       util.CommandFlags(flagsUpdate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "newName",
			Usage: "[optional] Update the server's name",
		},
		cli.StringFlag{
			Name:  "newIPv4",
			Usage: "[optional] Update the server's IPv4 address",
		},
		cli.StringFlag{
			Name:  "newIPv6",
			Usage: "[optional] Update the server's IPv6 address",
		},
	}
	return append(cf, idAndNameFlags...)
}

func commandUpdate(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	serverID := idOrName(c, client)
	opts := &osServers.UpdateOpts{
		Name:       c.String("newName"),
		AccessIPv4: c.String("newIPv4"),
		AccessIPv6: c.String("newIPv6"),
	}
	o, err := servers.Update(client, serverID, opts).Extract()
	if err != nil {
		fmt.Printf("Error updating server: %s\n", err)
		os.Exit(1)
	}
	output.Print(c, o, tableUpdate)
}

func tableUpdate(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{"property", "value"})
	keys := []string{"ID", "Name", "Public IPv4", "Public IPv6"}
	for _, key := range keys {
		tmp := ""
		switch key {
		case "Public IPv4":
			tmp = fmt.Sprint(m["AccessIPv4"])
		case "Public IPv6":
			tmp = fmt.Sprint(m["AccessIPv6"])
		default:
			tmp = fmt.Sprint(m[key])
		}
		if tmp == "<nil>" {
			tmp = ""
		}
		t.Append([]string{key, tmp})
	}
	t.Render()
}
