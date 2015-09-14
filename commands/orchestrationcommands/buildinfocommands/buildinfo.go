package buildinfocommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/buildinfo"
	"github.com/rackspace/rack/util"
)

var commandPrefix = "orchestration"

var BuildInfo = cli.Command{
	Name:        "buildinfo",
	Usage:       util.Usage(commandPrefix, "buildinfo", ""),
	Description: "Retrieve build information for a heat deployment",
	Action:      actionBuildInfo,
	Flags:       commandoptions.CommandFlags(flagsBuildInfo, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsBuildInfo, keysBuildInfo))
	},
}
var serviceClientType = "orchestration"

var keysBuildInfo = []string{"API", "Engine", "FusionAPI"}

type commandBuildInfo handler.Command

func flagsBuildInfo() []cli.Flag {
	return []cli.Flag{}
}

func actionBuildInfo(c *cli.Context) {
	command := &commandBuildInfo{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandBuildInfo) Context() *handler.Context {
	return command.Ctx
}

func (command *commandBuildInfo) Keys() []string {
	return keysBuildInfo
}

func (command *commandBuildInfo) ServiceClientType() string {
	return serviceClientType
}

func (command *commandBuildInfo) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandBuildInfo) Execute(resource *handler.Resource) {
	info, err := buildinfo.Get(command.Ctx.ServiceClient).Extract()
	if err != nil {
		resource.Err = err
	}
	result := structs.Map(info)
	for k, v := range result {
		result[k] = v.(map[string]interface{})["Revision"]
	}
	resource.Result = result
}
