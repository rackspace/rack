package instancecommands

import (
	"fmt"
	"time"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
	"github.com/rackspace/rack/util"
)

// delete is a reserved word in Go.
var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description: "Deletes an existing server",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
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
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the instance has been deleted",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	wait   bool
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
	wait := false
	if command.Ctx.CLIContext.IsSet("wait-for-completion") {
		wait = true
	}

	resource.Params = &paramsDelete{
		wait: wait,
	}
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

	if resource.Params.(*paramsDelete).wait {
		i := 0
		for i < 120 {
			_, err := servers.Get(command.Ctx.ServiceClient, serverID).Extract()
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
			i++
		}
		resource.Result = fmt.Sprintf("Deleted instance [%s]\n", serverID)
	} else {
		resource.Result = fmt.Sprintf("Deleting instance [%s]\n", serverID)
	}
}

func (command *commandDelete) StdinField() string {
	return "id"
}
