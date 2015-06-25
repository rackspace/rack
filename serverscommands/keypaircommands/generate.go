package keypaircommands

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var generate = cli.Command{
	Name:        "generate",
	Usage:       util.Usage(commandPrefix, "generate", "[--name <keypairName> | --stdin name]"),
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
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the keypair",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
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
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsGenerate).opts.Name = command.Ctx.CLIContext.String("name")
	return err
}

func (command *commandGenerate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsGenerate).opts
	keypair, err := keypairs.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = print(keypair)
}

func (command *commandGenerate) StdinField() string {
	return "name"
}

func print(kp *osKeypairs.KeyPair) string {
	output := []string{"PROPERTY\tVALUE",
		"Name\t\t%s",
		"Fingerprint\t%s",
		"PublicKey\t%s",
		"PrivateKey:\n%s",
	}
	return fmt.Sprintf(strings.Join(output, "\n"), kp.Name, kp.Fingerprint, kp.PublicKey, kp.PrivateKey)
}
