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

var rebuild = cli.Command{
	Name:        "rebuild",
	Usage:       util.Usage(commandPrefix, "rebuild", strings.Join([]string{util.IDOrNameUsage("instance"), "--image-id <image-id>", "--admin-pass <admin-pass>"}, " ")),
	Description: "Rebuilds an existing server",
	Action:      commandRebuild,
	Flags:       util.CommandFlags(flagsRebuild, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsRebuild, keysGet))
	},
}

func flagsRebuild() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "image-id",
			Usage: "[required] The ID of the image on which the server will be provisioned.",
		},
		cli.StringFlag{
			Name:  "admin-pass",
			Usage: "[required] The server's administrative password.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] The name for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "ipv4",
			Usage: "[optional] The IPv4 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "ipv6",
			Usage: "[optional] The IPv6 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string a key=value pairs.",
		},
	}
	return append(cf, util.IDAndNameFlags...)
}

func commandRebuild(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysGet,
	}
	err = util.CheckArgNum(c, 0)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	err = util.CheckFlagsSet(c, []string{"image-id", "admin-pass"})
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	imageID := c.String("image-id")

	opts := osServers.RebuildOpts{
		ImageID:    imageID,
		AdminPass:  c.String("admin-pass"),
		AccessIPv4: c.String("ipv4"),
		AccessIPv6: c.String("ipv6"),
	}

	if c.IsSet("metadata") {
		opts.Metadata, err = util.CheckKVFlag(c, "metadata")
		if err != nil {
			outputParams.Err = err
			output.Print(outputParams)
			return
		}
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

	if c.IsSet("rename") {
		opts.Name = c.String("rename")
	} else if c.IsSet("id") { // Must get the name from compute by ID
		serverResult, err := servers.Get(client, serverID).Extract()
		if err != nil {
			outputParams.Err = fmt.Errorf("Error retrieving instance [%s] for rebuild: %s\n", serverID, err)
			output.Print(outputParams)
			return
		}
		opts.Name = serverResult.Name
	} else if c.IsSet("name") {
		// Did not set rename, did not set id, can assume name
		opts.Name = c.String("name")
	}

	o, err := servers.Rebuild(client, serverID, opts).Extract()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error rebuilding instance [%s] with image [%s]: %s\n", serverID, imageID, err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		return serverSingle(o)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
