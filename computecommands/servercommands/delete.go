package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

// delete is a reserved word in Go.
var remove = cli.Command{
	Name:        "delete",
	Usage:       fmt.Sprintf("%s %s delete %s [optional flags]", util.Name, commandPrefix, idOrNameUsage),
	Description: "Deletes an existing server",
	Action:      commandDelete,
	Flags:       util.CommandFlags(flagsDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete))
	},
}

func flagsDelete() []cli.Flag {
	return idAndNameFlags
}

func commandDelete(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	serverID := idOrName(c, client)
	err := servers.Delete(client, serverID).ExtractErr()
	if err != nil {
		fmt.Printf("Error deleting server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
}
