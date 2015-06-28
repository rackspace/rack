package securitygroupcommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	securityGroups "github.com/rackspace/gophercloud/rackspace/networking/v2/security/groups"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", ""),
	Description: "Deletes an existing security group",
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
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the security group.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the security group.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysDelete = []string{"ID", "Name", "TenantID", "Rules"}

type paramsDelete struct {
	securityGroupID string
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
	resource.Params.(*paramsDelete).securityGroupID = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).securityGroupID = command.Ctx.CLIContext.String("id")
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	securityGroupID := resource.Params.(*paramsDelete).securityGroupID
	err := securityGroups.Delete(command.Ctx.ServiceClient, securityGroupID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted security group [%s]\n", securityGroupID)
}

func (command *commandDelete) StdinField() string {
	return "id"
}
