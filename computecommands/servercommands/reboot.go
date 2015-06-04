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

var reboot = cli.Command{
	Name:        "reboot",
	Usage:       fmt.Sprintf("%s %s reboot <serverID> [flags]", util.Name, commandPrefix),
	Description: "Reboots an existing server",
	Action:      commandReboot,
	Flags:       util.CommandFlags(flagsReboot),
}

func flagsReboot() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "method",
			Usage: `[required] The method to use for rebooting the server. Options are "hard" and "soft":
                "hard" will physically cut power to the machine and then restore it after a brief while;
                "soft" will ask the OS to restart under its own procedures.`,
		},
	}
}

func commandReboot(c *cli.Context) {
	util.CheckArgNum(c, 1)
	serverID := c.Args()[0]

	var how osServers.RebootMethod
	if !c.IsSet("method") {
		fmt.Printf("Required flag [method] for reboot not set.\n")
		os.Exit(1)
	}

	s := c.String("method")
	switch s {
	case "hard":
		how = osServers.PowerCycle
	case "soft":
		how = osServers.OSReboot
	default:
		fmt.Printf("Invalid option for reboot flag [method]: %s\n", s)
		os.Exit(1)
	}

	client := auth.NewClient("compute")
	err := servers.Reboot(client, serverID, how).ExtractErr()
	if err != nil {
		fmt.Printf("Error rebooting server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
}
