package stackcommands

import (
//    "github.com/rackspace/rack/commandoptions"
    "github.com/rackspace/rack/handler"
    "github.com/rackspace/rack/internal/github.com/codegangsta/cli"
    "github.com/rackspace/rack/internal/github.com/fatih/structs"
    osStacks "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
    "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stacks"
    "github.com/rackspace/rack/util"
    "fmt"
)

var list = cli.Command{
    Name: "list",
    Usage: util.Usage(commandPrefix, "list", ""),
    Description: "Retrieve a list of deployed stacks",
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

    pager := stacks.List(command.Ctx.ServiceClient, nil)
    //fmt.Println(pager)
    pages, err := pager.AllPages()
    if err != nil {
        resource.Err = err
        return
    }
    //fmt.Println(pages)
    info, err := osStacks.ExtractStacks(pages)
    //fmt.Println("Yeah!")
    if err != nil {
        resource.Err = err
        return
    }
    result := make([]map[string]interface{}, len(info))
    for j, stack := range info {
        result[j] = structs.Map(stack)
    }
    fmt.Println(result)
    resource.Result = result

}
