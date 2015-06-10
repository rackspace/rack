package instancecommands

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
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
	tmp := ""
	if b := util.ReadStdin(c); b != nil {
		tmp = string(b)
	} else {
		tmp = idOrName(c, client)
	}

	serverIDs := strings.Split(tmp, c.String("sep"))
	for _, serverID := range serverIDs {
		serverID := strings.TrimSpace(serverID)
		if serverID == "" {
			continue
		}
		fmt.Fprintf(c.App.Writer, "Deleting server: %s\n", serverID)
		err := servers.Delete(client, serverID).ExtractErr()
		if err != nil {
			fmt.Fprintf(c.App.Writer, "Error deleting server (%s): %s\n", serverID, err)
		}
	}
}
