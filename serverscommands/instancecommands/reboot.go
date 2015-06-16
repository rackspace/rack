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

var reboot = cli.Command{
	Name:        "reboot",
	Usage:       util.Usage(commandPrefix, "reboot", strings.Join([]string{util.IDOrNameUsage("instance"), "[--soft | --hard]"}, " ")),
	Description: "Reboots an existing server",
	Action:      commandReboot,
	Flags:       util.CommandFlags(flagsReboot, keysReboot),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsReboot, keysReboot))
	},
}

func flagsReboot() []cli.Flag {
	cf := []cli.Flag{
		cli.BoolFlag{
			Name:  "soft",
			Usage: "[optional; required if 'hard' is not provided] Ask the OS to restart under its own procedures.",
		},
		cli.BoolFlag{
			Name:  "hard",
			Usage: "[optional; required if 'soft' is not provided] Physically cut power to the machine and then restore it after a brief while.",
		},
	}
	return append(cf, util.IDAndNameFlags...)
}

var keysReboot = []string{}

func commandReboot(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysReboot,
	}
	err = util.CheckArgNum(c, 0)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	var how osServers.RebootMethod
	if c.IsSet("soft") {
		how = osServers.OSReboot
	}
	if c.IsSet("hard") {
		how = osServers.PowerCycle
	}

	if how == "" {
		outputParams.Err = util.Error(c, util.ErrMissingFlag{
			Msg: "One of either --soft or --hard must be provided.",
		})
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
	err = servers.Reboot(client, serverID, how).ExtractErr()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error retrieving server (%s): %s\n", serverID, err)
		output.Print(outputParams)
		return
	}
	f := func() interface{} {
		return fmt.Sprintf("Successfully rebooted instance [%s]", serverID)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
