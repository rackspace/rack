package instancecommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

// delete is a reserved word in Go.
var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", util.IDOrNameUsage("instance")),
	Description: "Deletes an existing server",
	Action:      commandDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return util.IDAndNameFlags
}

var keysDelete = []string{}

func commandDelete(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysDelete,
	}
	err = util.CheckArgNum(c, 0)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	outputParams.ServiceClientType = serviceClientType
	client, err := auth.NewClient(c, outputParams.ServiceClientType)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	serverID, err := util.IDOrName(c, client, osServers.IDFromName)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	err = servers.Delete(client, serverID).ExtractErr()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error deleting server (%s): %s\n", serverID, err)
		output.Print(outputParams)
		return
	}
	f := func() interface{} {
		return fmt.Sprintf("Successfully deleted instance [%s]", serverID)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
