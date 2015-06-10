package instancecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
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
	keys := []string{"ID", "Name", "Public IPv4", "Public IPv6"}
	mungeServerMap(m)
	output.MetaDataMapTable(c, m, keys)
}
