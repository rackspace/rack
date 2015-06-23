package keypaircommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--keypair <keypairName> | stdin keypair]"),
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
			Name:  "keypair",
			Usage: "[optional; required if `stdin` isn't provided] The name of the keypair",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `keypair` isn't provided] The field being piped into STDIN. Valid values are: keypair",
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
	err := command.Ctx.CheckFlagsSet([]string{"keypair"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).keypair = command.Ctx.CLIContext.String("keypair")
	return err
}

func (command *commandGet) Execute(resource *handler.Resource) {
	keypairName := resource.Params.(*paramsGet).keypair
	keypair, err := keypairs.Get(command.Ctx.ServiceClient, keypairName).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	result := structs.Map(keypair)
	if command.Ctx.CLIContext.GlobalIsSet("json") {
		resource.Result = result
	} else {
		// Assume they want the key directly
		resource.Result = result["PublicKey"]
	}
}

func (command *commandGet) StdinField() string {
	return "keypair"
}
