package servercommands

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
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
			Name:  "name",
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
		Name:       c.String("name"),
		AccessIPv4: c.String("accessIPv4"),
		AccessIPv6: c.String("accessIPv6"),
	}

	if c.IsSet("metadata") {
		metadata := make(map[string]string)
		metaStrings := strings.Split(c.String("metadata"), ",")
		for _, metaString := range metaStrings {
			temp := strings.Split(metaString, "=")
			if len(temp) != 2 {
				util.PrintError(c, util.ErrFlagFormatting{
					Msg: fmt.Sprintf("Expected key=value format but got %s for --metadata.\n", metaString),
				})
			}
			metadata[temp[0]] = temp[1]
		}
		opts.Metadata = metadata
	}

	client := auth.NewClient("compute")
	serverID := idOrName(c, client)
	o, err := servers.Rebuild(client, serverID, opts).Extract()
	if err != nil {
		fmt.Printf("Error rebuilding server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}
