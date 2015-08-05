package keypaircommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/rackspace/rack/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--name <keypairName> | --stdin name]"),
	Description: "Deletes a keypair",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the keypair",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	keypair string
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
	resource.Params.(*paramsDelete).keypair = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).keypair = command.Ctx.CLIContext.String("name")
	return err
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	keypairName := resource.Params.(*paramsDelete).keypair
	err := keypairs.Delete(command.Ctx.ServiceClient, keypairName).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted keypair [%s]\n", keypairName)
}

func (command *commandDelete) StdinField() string {
	return "name"
}
