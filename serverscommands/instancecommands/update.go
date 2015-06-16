package instancecommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", util.IDOrNameUsage("instance")),
	Description: "Updates an existing server",
	Action:      commandUpdate,
	Flags:       util.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "new-name",
			Usage: "[optional] Update the server's name",
		},
		cli.StringFlag{
			Name:  "new-ipv4",
			Usage: "[optional] Update the server's IPv4 address",
		},
		cli.StringFlag{
			Name:  "new-ipv6",
			Usage: "[optional] Update the server's IPv6 address",
		},
	}
	return append(cf, util.IDAndNameFlags...)
}

var keysUpdate = []string{"ID", "Name", "Public IPv4", "Public IPv6"}

func commandUpdate(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysUpdate,
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
	opts := &osServers.UpdateOpts{
		Name:       c.String("new-name"),
		AccessIPv4: c.String("new-ipv4"),
		AccessIPv6: c.String("new-ipv6"),
	}
	o, err := servers.Update(client, serverID, opts).Extract()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error updating instance [%s] : %s\n", serverID, err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		return serverSingle(o)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
