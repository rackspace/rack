package portcommands

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osPorts "github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/ports"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing ports",
	Action:      actionList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all ports. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list ports with this name.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "Only list ports that are up or not. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "network-id",
			Usage: "[optional] Only list ports with this network ID.",
		},
		cli.IntFlag{
			Name:  "status",
			Usage: "[optional] Only list ports that have this status.",
		},
		cli.StringFlag{
			Name:  "mac-address",
			Usage: "[optional] Only list ports with this MAC address.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] Only list ports that are owned by the tenant with this tenant ID.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing ports at this ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many ports at most.",
		},
	}
}

var keysList = []string{"ID", "Name", "Network ID", "Status", "MAC Address", "Device ID"}

type paramsList struct {
	opts     *osPorts.ListOpts
	allPages bool
}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &osPorts.ListOpts{
		Name:       c.String("name"),
		NetworkID:  c.String("network-id"),
		Status:     c.String("status"),
		MACAddress: c.String("mac-address"),
		TenantID:   c.String("tenant-id"),
		Marker:     c.String("marker"),
		Limit:      c.Int("limit"),
	}

	if c.IsSet("up") {
		upRaw := c.String("up")
		up, err := strconv.ParseBool(upRaw)
		if err != nil {
			return fmt.Errorf("Invalid value for flag `up`: %s. Options are: true, false", upRaw)
		}
		opts.AdminStateUp = &up
	}

	resource.Params = &paramsList{
		opts:     opts,
		allPages: c.Bool("all-pages"),
	}

	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	allPages := resource.Params.(*paramsList).allPages
	pager := ports.List(command.Ctx.ServiceClient, opts)
	if allPages {
		pages, err := pager.AllPages()
		if err != nil {
			resource.Err = err
			return
		}
		info, err := osPorts.ExtractPorts(pages)
		if err != nil {
			resource.Err = err
			return
		}
		result := make([]map[string]interface{}, len(info))
		for j, subnet := range info {
			result[j] = portSingle(&subnet)
		}
		resource.Result = result
	} else {
		limit := opts.Limit
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			info, err := osPorts.ExtractPorts(page)
			if err != nil {
				return false, err
			}
			result := make([]map[string]interface{}, len(info))
			for j, subnet := range info {
				result[j] = portSingle(&subnet)
			}
			resource.Result = result
			if len(info) >= limit {
				return false, nil
			}
			limit -= len(info)
			command.Ctx.WaitGroup.Add(1)
			command.Ctx.Results <- resource
			return true, nil
		})
		if err != nil {
			resource.Err = err
			return
		}
	}
}
