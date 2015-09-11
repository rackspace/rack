package stackresourcecommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackresources"
	"github.com/rackspace/rack/util"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
)

var getMetadata = cli.Command{
	Name:        "get-metadata",
	Usage:       util.Usage(commandPrefix, "get-metadata", "--name <stackName> --id <stackID> --resource <resourceName>"),
	Description: "Show metadata for specified resource",
	Action:      actionGetMetadata,
	Flags:       commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata))
	},
}

func flagsGetMetadata() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The stack id.",
		},
		cli.StringFlag{
			Name:  "resource",
			Usage: "[required] The resource name.",
		},
	}
}

type paramsGetMetadata struct {
	stackName    string
	stackID      string
	resourceName string
}

var keysGetMetadata = []string{"Name", "PhysicalID", "Type", "Status", "UpdatedTime"}

type commandGetMetadata handler.Command

func actionGetMetadata(c *cli.Context) {
	command := &commandGetMetadata{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGetMetadata) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGetMetadata) Keys() []string {
	return keysGetMetadata
}

func (command *commandGetMetadata) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGetMetadata) HandleFlags(resource *handler.Resource) error {
	if err := command.Ctx.CheckFlagsSet([]string{"resource"}); err != nil {
		return err
	}
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}

	resource.Params = &paramsGetMetadata{
		stackName:    name,
		stackID:      id,
		resourceName: command.Ctx.CLIContext.String("resource"),
	}
	return nil
}

func (command *commandGetMetadata) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGetMetadata)
	stackName := params.stackName
	stackID := params.stackID
	resourceName := params.resourceName
	metadata, err := stackresources.Metadata(command.Ctx.ServiceClient, stackName, stackID, resourceName).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	if len(metadata) == 0 {
		resource.Result = "None"
	} else {
		resource.Result = metadata
	}
}
