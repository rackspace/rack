package stackcommands

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var abandon = cli.Command{
	Name:        "abandon",
	Usage:       util.Usage(commandPrefix, "abandon", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description: "Deletes an existing stack, but leave resources intact",
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
		cli.StringFlag{
			Name:  "output-file",
			Usage: "[optional] The file into which result of abandon is stored.",
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
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}

	resource.Params = &paramsAbandon{
		stackName: name,
		stackID:  id,
	}
	return nil
}

func (command *commandAbandon) Execute(resource *handler.Resource) {
	c := command.Ctx.CLIContext
	params := resource.Params.(*paramsAbandon)
	stackName := params.stackName
	stackID := params.stackID
	res := stacks.Abandon(command.Ctx.ServiceClient, stackName, stackID)
	res_stack, err := res.Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = stackSingle(res_stack)
	if c.IsSet("output-file") {
		data, err := json.MarshalIndent(res.Body, "", "    ")
		if err != nil {
			resource.Err = err
			return
		}
		abs, err := filepath.Abs(c.String("output-file"))
		if err != nil {
			resource.Err = err
			return
		}
		if err := ioutil.WriteFile(abs, data, 0644); err != nil {
			resource.Err = err
			return
		}
	}
}
