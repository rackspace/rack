package stackcommands

import (
	"encoding/json"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var abandon = cli.Command{
	Name:        "abandon",
	Usage:       util.Usage(commandPrefix, "abandon", "[--name <stackName> | --id <stackID> | --stdin name]"),
	Description: "Deletes an existing stack, but leaves resources intact",
	Action:      actionAbandon,
	Flags:       commandoptions.CommandFlags(flagsAbandon, keysAbandon),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsAbandon, keysAbandon))
	},
}

func flagsAbandon() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the stack.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the stack.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
	}
}

var keysAbandon = []string{"Status", "Name", "Template", "Action", "ID", "Resources", "Files", "StackUserProjectID", "ProjectID", "Environment"}

type paramsAbandon struct {
	stackName string
	stackID   string
}

type commandAbandon handler.Command

func actionAbandon(c *cli.Context) {
	command := &commandAbandon{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandAbandon) Context() *handler.Context {
	return command.Ctx
}

func (command *commandAbandon) Keys() []string {
	return keysAbandon
}

func (command *commandAbandon) ServiceClientType() string {
	return serviceClientType
}

func (command *commandAbandon) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandAbandon) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params.(*paramsAbandon).stackName = name
	resource.Params.(*paramsAbandon).stackID = id
	return nil
}

func (command *commandAbandon) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsAbandon{
		stackName: name,
		stackID:   id,
	}

	resource.Params.(*paramsAbandon).stackName = name
	resource.Params.(*paramsAbandon).stackID = id
	return nil
}

func (command *commandAbandon) Execute(resource *handler.Resource) {

	params := resource.Params.(*paramsAbandon)
	stackName := params.stackName
	stackID := params.stackID
	res := stacks.Abandon(command.Ctx.ServiceClient, stackName, stackID)
	resStack, err := res.Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = resStack
}

func (command *commandAbandon) StdinField() string {
	return "name"
}

func (command *commandAbandon) PreTable(resource *handler.Resource) error {
	resource.Result = stackSingle(resource.Result)
	return nil
}

func (command *commandAbandon) PreJSON(resource *handler.Resource) error {
	var resInterface map[string]interface{}
	res, err := json.Marshal(resource.Result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &resInterface)
	if err != nil {
		return err
	}
	resource.Result = resInterface
	resource.Keys = []string{"status", "name", "template", "action", "id", "resources", "files", "stack_user_project_id", "project_id", "environment"}
	return nil
}

func (command *commandAbandon) PreCSV(resource *handler.Resource) error {
	command.PreTable(resource)
	return nil
}
