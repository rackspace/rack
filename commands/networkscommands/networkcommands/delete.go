package networkcommands

import (
	"fmt"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osNetworks "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/networking/v2/networks"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", ""),
	Description: "Deletes an existing network",
	Action:      actionDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the network",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the network.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	networkID string
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
	resource.Params.(*paramsDelete).networkID = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	networkID, err := command.Ctx.IDOrName(osNetworks.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).networkID = networkID
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	networkID := resource.Params.(*paramsDelete).networkID
	err := networks.Delete(command.Ctx.ServiceClient, networkID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted network [%s] \n", networkID)
}

func (command *commandDelete) StdinField() string {
	return "id"
}
