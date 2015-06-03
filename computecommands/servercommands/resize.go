package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var resize = cli.Command{
	Name:        "resize",
	Usage:       fmt.Sprintf("%s %s resize <serverID> [flags]", util.Name, commandPrefix),
	Description: "Rebuilds an existing server",
	Action:      commandResize,
	Flags:       flagsResize(),
}

func flagsResize() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "flavorID",
			Usage: "[required] The ID of the flavor that the resized server should have.",
		},
	}
}

func commandResize(c *cli.Context) {
	util.CheckArgNum(c, 1)
	serverID := c.Args()[0]

	if !c.IsSet("flavorID") {
		fmt.Printf("Required flag [flavorID] for resize not set.\n")
		os.Exit(1)
	}

	opts := osServers.ResizeOpts{
		FlavorRef: c.String("flavorID"),
	}

	client := auth.NewClient("compute")
	err := servers.Resize(client, serverID, opts).ExtractErr()
	if err != nil {
		fmt.Printf("Error resizing server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
}
