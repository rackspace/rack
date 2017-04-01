package instancecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing database instances",
	Action:      actionList,
	Flags:       commandoptions.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "include-replicas",
			Usage: "[optional] Include replica instances. Default is FALSE.",
		},
		cli.BoolFlag{
			Name:  "include-ha",
			Usage: "[optional] Include High Availability instances. Default is FALSE.",
		},
	}
}

var keysList = []string{"ID", "Name", "Flavor", "Size", "Datastore", "Status"}

type paramsList struct {
	opts     *instances.ListOpts
	allPages bool
}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	resource.Params = &paramsList{
		opts: &instances.ListOpts{
			IncludeReplicas: c.Bool("include-replicas"),
			IncludeHA:       c.Bool("include-ha"),
		},
	}
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsList).opts
	pager := instances.List(command.Ctx.ServiceClient, opts)

	pages, err := pager.AllPages()
	if err != nil {
		resource.Err = err
		return
	}

	info, err := instances.ExtractInstances(pages)
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("%#v", info)

	result := make([]map[string]interface{}, len(info))
	for j, instance := range info {
		result[j] = singleInstance(&instance)
	}

	resource.Result = result
}
