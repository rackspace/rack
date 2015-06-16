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

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing servers",
	Action:      commandList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list servers with this name.",
		},
		cli.StringFlag{
			Name:  "changes-since",
			Usage: "Only list servers that have been changed since this time/date stamp.",
		},
		cli.StringFlag{
			Name:  "image",
			Usage: "Only list servers that have this image ID.",
		},
		cli.StringFlag{
			Name:  "flavor",
			Usage: "Only list servers that have this flavor ID.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list servers that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "Start listing servers at this server ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "Only return this many servers at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Status", "Public IPv4", "Private IPv4", "Image", "Flavor"}

func commandList(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysList,
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

	opts := osServers.ListOpts{
		ChangesSince: c.String("changes-since"),
		Image:        c.String("image"),
		Flavor:       c.String("flavor"),
		Name:         c.String("name"),
		Status:       c.String("status"),
		Marker:       c.String("marker"),
		Limit:        c.Int("limit"),
	}
	allPages, err := servers.List(client, opts).AllPages()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error listing servers: %s\n", err)
		output.Print(outputParams)
		return
	}
	o, err := servers.ExtractServers(allPages)
	if err != nil {
		outputParams.Err = fmt.Errorf("Error listing servers: %s\n", err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, server := range o {
			m[j] = serverSingle(&server)
		}
		return m
	}
	outputParams.F = &f
	output.Print(outputParams)
}
