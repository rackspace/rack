package instancecommands

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var resize = cli.Command{
	Name:        "resize",
	Usage:       util.Usage(commandPrefix, "resize", strings.Join([]string{util.IDOrNameUsage("instance"), "--flavor-id <flavor-id>"}, " ")),
	Description: "Resizes an existing server",
	Action:      actionResize,
	Flags:       util.CommandFlags(flagsResize, keysResize),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsResize, keysResize))
	},
}

func flagsResize() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "flavor-id",
			Usage: "[required] The ID of the flavor that the resized server should have.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
	return append(cf, util.IDAndNameFlags...)
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
	resource.Result = fmt.Sprintf("Successfully resized instance [%s] to flavor [%s]", params.serverID, params.opts.FlavorRef)
}

func (command *commandResize) StdinField() string {
	return "id"
}
