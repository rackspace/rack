package instancecommands

import (
	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", util.IDOrNameUsage("instance")),
	Description: "Updates an existing server",
	Action:      actionUpdate,
	Flags:       util.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "new-name",
			Usage: "[optional] Update the server's name",
		},
		cli.StringFlag{
			Name:  "new-ipv4",
			Usage: "[optional] Update the server's IPv4 address",
		},
		cli.StringFlag{
			Name:  "new-ipv6",
			Usage: "[optional] Update the server's IPv6 address",
		},
	}
	return append(cf, util.IDAndNameFlags...)
}

var keysUpdate = []string{"ID", "Name", "Public IPv4", "Public IPv6"}

type paramsUpdate struct {
	serverID string
	opts     *osServers.UpdateOpts
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
	opts := &osServers.UpdateOpts{
		Name:       c.String("new-name"),
		AccessIPv4: c.String("new-ipv4"),
		AccessIPv6: c.String("new-ipv6"),
	}
	resource.Params = &paramsUpdate{
		opts: opts,
	}
	return nil
}

func (command *commandUpdate) HandleSingle(resource *handler.Resource) error {
	id, err := command.Ctx.IDOrName(osServers.IDFromName)
	resource.Params.(*paramsUpdate).serverID = id
	return err
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpdate)
	server, err := servers.Update(command.Ctx.ServiceClient, params.serverID, params.opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = serverSingle(server)
}
