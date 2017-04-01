package instancecommands

import (
	"fmt"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
	"github.com/rackspace/rack/util"
)

var getConfig = cli.Command{
	Name:        "get-config",
	Usage:       util.Usage(commandPrefix, "get-config", "[--id <serverID> | --stdin id]"),
	Description: "Lists the default configuration settings for an existing database instance",
	Action:      actionGetConfig,
	Flags:       commandoptions.CommandFlags(flagsGetConfig, keysGetConfig),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetConfig, keysGetConfig))
	},
}

func flagsGetConfig() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the database instance.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
	}
}

var keysGetConfig = []string{}

type paramsGetConfig struct {
	id string
}

type commandGetConfig handler.Command

func actionGetConfig(c *cli.Context) {
	command := &commandGetConfig{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGetConfig) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGetConfig) Keys() []string {
	return keysGetConfig
}

func (command *commandGetConfig) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGetConfig) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGetConfig{}
	return nil
}

func (command *commandGetConfig) StdinField() string {
	return "id"
}

func (command *commandGetConfig) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGetConfig).id = item
	return nil
}

func (command *commandGetConfig) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGetConfig).id = command.Ctx.CLIContext.String("id")
	return err
}

func (command *commandGetConfig) Execute(resource *handler.Resource) {
	id := resource.Params.(*paramsGetConfig).id
	config, err := instances.GetDefaultConfig(command.Context().ServiceClient, id).Extract()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = config
}

func (command *commandGetConfig) PreCSV(resource *handler.Resource) error {
	key := "Config"
	resource.Result = map[string]interface{}{
		key: resource.Result,
	}
	resource.Keys = []string{key}

	keys := resource.Keys
	res := resource.Result.(map[string]interface{})
	if m, ok := res[key]; ok && util.Contains(keys, key) {
		for k, v := range m.(map[string]string) {
			newKey := fmt.Sprintf("%s", k)
			res[newKey] = v
			keys = append(keys, newKey)
		}
	}
	delete(res, key)
	resource.Keys = util.RemoveFromList(keys, key)

	return nil
}

func (command *commandGetConfig) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
