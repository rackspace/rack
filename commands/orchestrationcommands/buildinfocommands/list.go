package buildinfocommands

import (
    "github.com/rackspace/rack/handler"
    "github.com/rackspace/rack/internal/github.com/codegangsta/cli"
    "github.com/rackspace/rack/internal/github.com/fatih/structs"
    "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/buildinfo"
    "github.com/rackspace/rack/util"
)

var list = cli.Command{
    Name: "list",
    Usage: util.Usage(commandPrefix, "list", ""),
    Description: "Retrieve build information for a heat deployment",
    Action: actionGet,
}

var keysList = []string{"API", "Engine", "FusionAPI"}

type commandGet handler.Command

func actionGet(c *cli.Context) {
    command := &commandGet {
        Ctx: &handler.Context{
            CLIContext: c,
        },
    }
    handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
    return command.Ctx
}

func (command *commandGet) Keys() []string {
    return keysList
}

func (command *commandGet) ServiceClientType() string {
    return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
    return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
    buildids, err := buildinfo.Get(command.Ctx.ServiceClient).Extract()
    if err != nil {
        resource.Err = err
        return
    }
    resource.Result = structs.Map(buildids)
}
