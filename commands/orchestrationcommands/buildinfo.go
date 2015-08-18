package orchestrationcommands

import (
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/buildinfo"
)

var buildInfo = cli.Command{
	Name:        "buildinfo",
	Usage:       "buildinfo",
	Description: "Retrieve build information for a heat deployment",
	Action:      actionBuildInfo,
}

var keysBuildInfo = []string{"API", "Engine", "FusionAPI"}

type commandBuildInfo handler.Command

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
	result := buildinfo.Get(command.Ctx.ServiceClient)
	resource.Result = result.PrettyPrintJSON()
}
