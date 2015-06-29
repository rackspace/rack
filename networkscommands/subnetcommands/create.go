package subnetcommands

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osSubnets "github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/subnets"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--network-id <networkID> --cidr <CIDR> --ip-version <4|6>"),
	Description: "Creates a subnet",
	Action:      actionCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "network-id",
			Usage: "[required] The network ID under which to create this subnet.",
		},
		cli.StringFlag{
			Name:  "cidr",
			Usage: "[required] The CIDR of this subnet.",
		},
		cli.IntFlag{
			Name:  "ip-version",
			Usage: "[required] The IP version this subnet should have. Options are: 4, 6.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] The name this subnet should have.",
		},
		cli.StringFlag{
			Name:  "gateway-ip",
			Usage: "[optional] The gateway IP address this subnet should have.",
		},
		cli.BoolFlag{
			Name:  "enable-dhcp",
			Usage: "[optional] If set, DHCP will be enabled on this subnet.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant that should own this subnet.",
		},
		cli.StringSliceFlag{
			Name: "allocation-pool",
			Usage: strings.Join([]string{"[optional] An allocation pool for this subnet. This flag may be provided several times.\n",
				"\tEach one of these flags takes 2 values: start and end.\n",
				"\tExamle: --allocation-pool start=192.0.2.1,end=192.0.2.254 --allocation-pool start:172.20.0.1,end=172.20.0.254"}, ""),
		},
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

var keysCreate = []string{"ID", "Name", "Network ID", "CIDR", "EnableDHCP", "Gateway IP", "DNS Nameservers", "Allocation Pools"}

type paramsCreate struct {
	opts *osSubnets.CreateOpts
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"network-id", "cidr", "ip-version"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := &osSubnets.CreateOpts{
		NetworkID: c.String("network-id"),
		CIDR:      c.String("cidr"),
		Name:      c.String("name"),
		TenantID:  c.String("tenant-id"),
		GatewayIP: c.String("gateway-ip"),
		IPVersion: c.Int("ip-version"),
	}

	if c.IsSet("enable-dhcp") {
		enableDHCP := true
		opts.EnableDHCP = &enableDHCP
	}

	if c.IsSet("dns-nameservers") {
		opts.DNSNameservers = strings.Split(c.String("dns-nameservers"), ",")
	}

	if c.IsSet("allocation-pool") {
		allocationPoolsRaw := c.StringSlice("allocation-pool")
		allocationPoolsRawSlice, err := command.Ctx.CheckStructFlag(allocationPoolsRaw)
		if err != nil {
			return err
		}
		allocationPools := make([]osSubnets.AllocationPool, len(allocationPoolsRawSlice))
		for i, allocationPoolMap := range allocationPoolsRawSlice {
			allocationPools[i] = osSubnets.AllocationPool{
				Start: allocationPoolMap["start"].(string),
				End:   allocationPoolMap["end"].(string),
			}
		}
		opts.AllocationPools = allocationPools
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

	resource.Params = &paramsCreate{
		opts: opts,
	}
	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	subnet, err := subnets.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = subnetSingle(subnet)
}
