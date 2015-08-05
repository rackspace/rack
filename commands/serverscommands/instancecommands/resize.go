package instancecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osServers "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

var resize = cli.Command{
	Name:        "resize",
	Usage:       util.Usage(commandPrefix, "resize", "[--id <serverID>|--name <serverName>|--stdin id] --flavor-id <flavor-id>"),
	Description: "Resizes an existing server",
	Action:      actionResize,
	Flags:       commandoptions.CommandFlags(flagsResize, keysResize),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsResize, keysResize))
	},
}

func flagsResize() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "flavor-id",
			Usage: "[required] The ID of the flavor that the resized server should have.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the server.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the server.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysResize = []string{}

type paramsResize struct {
	serverID string
	opts     *osServers.ResizeOpts
}

type commandResize handler.Command

func actionResize(c *cli.Context) {
	command := &commandResize{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandResize) Context() *handler.Context {
	return command.Ctx
}

func (command *commandResize) Keys() []string {
	return keysResize
}

func (command *commandResize) ServiceClientType() string {
	return serviceClientType
}

func (command *commandResize) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"flavor-id"})
	if err != nil {
		return err
	}
	flavorID := command.Ctx.CLIContext.String("flavor-id")
	opts := &osServers.ResizeOpts{
		FlavorRef: flavorID,
	}
	resource.Params = &paramsResize{
		opts: opts,
	}
	return nil
}

func (command *commandResize) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsResize).serverID = item
	return nil
}

func (command *commandResize) HandleSingle(resource *handler.Resource) error {
	id, err := command.Ctx.IDOrName(osServers.IDFromName)
	resource.Params.(*paramsResize).serverID = id
	return err
}

func (command *commandResize) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsResize)
	err := servers.Resize(command.Ctx.ServiceClient, params.serverID, params.opts).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Successfully resized instance [%s] to flavor [%s]\n", params.serverID, params.opts.FlavorRef)
}

func (command *commandResize) StdinField() string {
	return "id"
}
