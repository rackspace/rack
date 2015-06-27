package subnetcommands

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSubnets "github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/subnets"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", ""),
	Description: "Updates a subnet",
	Action:      actionUpdate,
	Flags:       util.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The ID of the subnet to update.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The name of the subnet to update.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] The name to which this subnet should be updated.",
		},
		cli.StringFlag{
			Name:  "gateway-ip",
			Usage: "[optional] The gateway IP address this subnet should have.",
		},
		/*
			cli.StringFlag{
				Name:  "enable-dhcp",
				Usage: "[optional] Whether or not DHCP should be enabled on this subnet. Options are: true, false",
			},
		*/
		cli.StringFlag{
			Name:  "dns-nameservers",
			Usage: "[optional] A comma-separated list of DNS Nameservers for this subnet.",
		},
		/*
			cli.StringSliceFlag{
				Name: "host-route",
				Usage: strings.Join([]string{"[optional] A host route for this subnet. This flag may be provided several times.\n",
					"\tEach one of these flags takes 2 values: dest (the destination CIDR) and next (the next hop).\n",
					"\tExamle: --host-route dest=40.0.1.0/24,next=40.0.0.2"}, ""),
			},
		*/
	}
}

var keysUpdate = []string{"ID", "Name", "Network ID", "CIDR", "EnableDHCP", "Gateway IP", "DNS Nameservers", "Allocation Pools"}

type paramsUpdate struct {
	subnetID string
	opts     *osSubnets.UpdateOpts
}

type commandUpdate handler.Command

func actionUpdate(c *cli.Context) {
	command := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpdate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpdate) Keys() []string {
	return keysUpdate
}

func (command *commandUpdate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpdate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext

	opts := &osSubnets.UpdateOpts{
		Name:      c.String("rename"),
		GatewayIP: c.String("gateway-ip"),
	}

	/*
		if c.IsSet("enable-dhcp") {
			enableDHCPRaw := c.String("enable-dhcp")
			enableDHCP, err := strconv.ParseBool(enableDHCPRaw)
			if err != nil {
				return fmt.Errorf("Invalid value for flag `shared`: %s. Options are: true, false", enableDHCPRaw)
			}
			opts.EnableDHCP = &enableDHCP
		}
	*/

	if c.IsSet("dns-nameservers") {
		opts.DNSNameservers = strings.Split(c.String("dns-nameservers"), ",")
	}

	/*
		if c.IsSet("host-route") {
			hostRoutesRaw := c.StringSlice("host-route")
			hostRoutesRawSlice, err := command.Ctx.CheckStructFlag(hostRoutesRaw)
			if err != nil {
				return err
			}
			hostRoutes := make([]osSubnets.HostRoute, len(hostRoutesRawSlice))
			for i, hostRouteMap := range hostRoutesRawSlice {
				hostRoutes[i] = osSubnets.HostRoute{
					DestinationCIDR: hostRouteMap["dest"].(string),
					NextHop:         hostRouteMap["next"].(string),
				}
			}
			opts.HostRoutes = hostRoutes
		}
	*/

	resource.Params = &paramsUpdate{
		opts: opts,
	}

	return nil
}

func (command *commandUpdate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsUpdate).subnetID = command.Ctx.CLIContext.String("id")
	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	subnetID := resource.Params.(*paramsUpdate).subnetID
	opts := resource.Params.(*paramsUpdate).opts
	subnet, err := subnets.Update(command.Ctx.ServiceClient, subnetID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = subnetSingle(subnet)
}
