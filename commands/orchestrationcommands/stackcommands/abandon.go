package stackcommands

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
	"github.com/rackspace/rack/util"
)

var abandon = cli.Command{
	Name:        "abandon",
	Usage:       util.Usage(commandPrefix, "abandon", "[--name <stackName> | --id <stackID> | --stdin name]"),
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
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
		cli.StringFlag{
			Name:  "output-file",
			Usage: "[optional] The file into which result of abandon is stored.",
		},
	}
}

var keysAbandon = []string{"Status", "Name", "Template", "Action", "ID", "Resources", "Files", "StackUserProjectID", "ProjectID", "Environment"}

type paramsAbandon struct {
	stackName  string
	stackID    string
	outputFile string
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
	var outputFile string = ""
	if c.IsSet("output-file") {
		outputFile = c.String("output-file")
	}
	resource.Params = &paramsAbandon{
		outputFile: outputFile,
	}
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
	resource.Params.(*paramsAbandon).stackName = name
	resource.Params.(*paramsAbandon).stackID = id
	return nil
}

func (command *commandAbandon) Execute(resource *handler.Resource) {

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
	if params.outputFile != "" {
		data, err := json.MarshalIndent(res.Body, "", "    ")
		if err != nil {
			resource.Err = err
			return
		}
		abs, err := filepath.Abs(params.outputFile)
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

func (command *commandAbandon) StdinField() string {
	return "name"
}
