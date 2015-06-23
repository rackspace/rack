package keypaircommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var generate = cli.Command{
	Name:        "generate",
	Usage:       util.Usage(commandPrefix, "generate", "[--keypair <keypairName> | --stdin keypair]"),
	Description: "Generates a keypair",
	Action:      actionGenerate,
	Flags:       util.CommandFlags(flagsGenerate, keysGenerate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGenerate, keysGenerate))
	},
}

func flagsGenerate() []cli.Flag {
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

var keysGenerate = []string{"Name", "Fingerprint", "PublicKey", "PrivateKey"}

type paramsGenerate struct {
	opts *osKeypairs.CreateOpts
}

type commandGenerate handler.Command

func actionGenerate(c *cli.Context) {
	command := &commandGenerate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGenerate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGenerate) Keys() []string {
	return keysGenerate
}

func (command *commandGenerate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGenerate) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGenerate{
		opts: &osKeypairs.CreateOpts{},
	}
	return nil
}

func (command *commandGenerate) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsGenerate).opts.Name = item
	return nil
}

func (command *commandGenerate) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"keypair"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGenerate).opts.Name = command.Ctx.CLIContext.String("keypair")
	return err
}

func (command *commandGenerate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsGenerate).opts
	keypair, err := keypairs.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(keypair)
}

func (command *commandGenerate) StdinField() string {
	return "keypair"
}
