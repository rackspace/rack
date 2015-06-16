package instancecommands

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var resize = cli.Command{
	Name:        "resize",
	Usage:       util.Usage(commandPrefix, "resize", strings.Join([]string{util.IDOrNameUsage("instance"), "--flavor-id <flavor-id>"}, " ")),
	Description: "Resizes an existing server",
	Action:      commandResize,
	Flags:       util.CommandFlags(flagsResize, keysResize),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsResize, keysResize))
	},
}

func flagsResize() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "flavor-id",
			Usage: "[required] The ID of the flavor that the resized server should have.",
		},
	}
	return append(cf, util.IDAndNameFlags...)
}

var keysResize = []string{}

func commandResize(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysResize,
	}
	err = util.CheckArgNum(c, 0)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	err = util.CheckFlagsSet(c, []string{"flavor-id"})
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

	flavorID := c.String("flavor-id")
	opts := osServers.ResizeOpts{
		FlavorRef: flavorID,
	}
	err = servers.Resize(client, serverID, opts).ExtractErr()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error resizing instance [%s] to flavor [%s]: %s\n", serverID, flavorID, err)
		output.Print(outputParams)
		return
	}
	f := func() interface{} {
		return fmt.Sprintf("Successfully resized instance [%s] to flavor [%s]", serverID, flavorID)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
