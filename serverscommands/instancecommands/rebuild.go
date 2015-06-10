package instancecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var rebuild = cli.Command{
	Name:        "rebuild",
	Usage:       fmt.Sprintf("%s %s rebuild %s [--imageID <imageID>] [--adminPass <adminPass>] [optional flags]", util.Name, commandPrefix, idOrNameUsage),
	Description: "Rebuilds an existing server",
	Action:      commandRebuild,
	Flags:       util.CommandFlags(flagsRebuild),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsRebuild))
	},
}

func flagsRebuild() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "imageID",
			Usage: "[required] The ID of the image on which the server will be provisioned.",
		},
		cli.StringFlag{
			Name:  "adminPass",
			Usage: "[required] The server's administrative password.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] The name for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "accessIPv4",
			Usage: "[optional] The IPv4 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "accessIPv6",
			Usage: "[optional] The IPv6 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string a key=value pairs.",
		},
	}
	return append(cf, idAndNameFlags...)
}

func commandRebuild(c *cli.Context) {
	util.CheckArgNum(c, 0)

	if !c.IsSet("imageID") {
		util.PrintError(c, util.ErrMissingFlag{
			Msg: "--imageID is required.",
		})
	}

	if !c.IsSet("adminPass") {
		util.PrintError(c, util.ErrMissingFlag{
			Msg: "--adminPass is required.",
		})
	}

	opts := osServers.RebuildOpts{
		ImageID:    c.String("imageID"),
		AdminPass:  c.String("adminPass"),
		AccessIPv4: c.String("accessIPv4"),
		AccessIPv6: c.String("accessIPv6"),
	}

	if c.IsSet("rename") {
		opts.Name = c.String("rename")
	}

	if c.IsSet("metadata") {
		opts.Metadata = util.CheckKVFlag(c, "metadata")
	}

	client := auth.NewClient("compute")
	serverID := idOrName(c, client)

	if c.IsSet("rename") {
		opts.Name = c.String("rename")
	} else if c.IsSet("id") { // Must get the name from compute by ID
		getResult := osServers.Get(client, serverID)
		serverResult, err := getResult.Extract()
		if err != nil {
			fmt.Printf("Error rebuilding server (%s): %s\n", serverID, err)
			os.Exit(1)
		}
		opts.Name = serverResult.Name
	} else if c.IsSet("name") {
		// Did not set rename, did not set id, can assume name
		opts.Name = c.String("name")
	}

	o, err := servers.Rebuild(client, serverID, opts).Extract()
	if err != nil {
		fmt.Printf("Error rebuilding server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}
