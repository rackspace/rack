package instancecommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

// delete is a reserved word in Go.
var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", util.IDOrNameUsage("instance")),
	Description: "Deletes an existing server",
	Action:      actionDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	cf := []cli.Flag{
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
	return append(cf, util.IDAndNameFlags...)
}

var keysDelete = []string{}

type paramsDelete struct {
	server string
}

type commandDelete handler.Command

func actionDelete(c *cli.Context) {
	command := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDelete) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDelete) Keys() []string {
	return keysDelete
}

func (command *commandDelete) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDelete) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsDelete{}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).server = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	id, err := command.Ctx.IDOrName(osServers.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).server = id
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	serverID := resource.Params.(*paramsDelete).server
	err := servers.Delete(command.Ctx.ServiceClient, serverID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Deleting instance [%s]", serverID)
}

func (command *commandDelete) StdinField() string {
	return "id"
}
