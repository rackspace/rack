package keypaircommands

import (
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/jrperritt/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--name <keypairName> | stdin name]"),
	Description: "Retreives a keypair",
	Action:      actionGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the keypair",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

var keysGet = []string{"Name", "Fingerprint", "PublicKey", "UserID"}

type paramsGet struct {
	keypair string
}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
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
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGet).keypair = item
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).keypair = command.Ctx.CLIContext.String("name")
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	keypairName := resource.Params.(*paramsGet).keypair
	keypair, err := keypairs.Get(command.Ctx.ServiceClient, keypairName).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(keypair)
}

func (command *commandGet) StdinField() string {
	return "name"
}

func (command *commandGet) CSV(resource *handler.Resource) {
	resource.Result = resource.Result.(map[string]interface{})["PublicKey"]
}

func (command *commandGet) Table(resource *handler.Resource) {
	command.CSV(resource)
}
