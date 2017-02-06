package keypaircommands

import (
	"fmt"
	"io/ioutil"

	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/handler"
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/rackspace/rack/util"
)

var upload = cli.Command{
	Name:        "upload",
	Usage:       util.Usage(commandPrefix, "upload", "--name <keypairName> [--public-key <publicKey> | --file <file>]"),
	Description: "Uploads a keypair",
	Action:      actionUpload,
	Flags:       commandoptions.CommandFlags(flagsUpload, keysUpload),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpload, keysUpload))
	},
}

func flagsUpload() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the keypair",
		},
		cli.StringFlag{
			Name:  "public-key",
			Usage: "[optional; required if `file` is not provided] The public ssh key to associate with the user's account.",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "[optional; required if `public-key` is not provided] The name of the file containing the public key.",
		},
	}
}

var keysUpload = []string{"Name", "Fingerprint", "PublicKey", "PrivateKey"}

type paramsUpload struct {
	opts *osKeypairs.CreateOpts
}

type commandUpload handler.Command

func actionUpload(c *cli.Context) {
	command := &commandUpload{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpload) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpload) Keys() []string {
	return keysUpload
}

func (command *commandUpload) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpload) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}

	opts := &osKeypairs.CreateOpts{
		Name: command.Ctx.CLIContext.String("name"),
	}

	if command.Ctx.CLIContext.IsSet("file") {
		s := command.Ctx.CLIContext.String("file")
		pk, err := ioutil.ReadFile(s)
		if err != nil {
			return err
		}
		opts.PublicKey = string(pk)
	} else if command.Ctx.CLIContext.IsSet("public-key") {
		s := command.Ctx.CLIContext.String("public-key")
		opts.PublicKey = s
	} else {
		return fmt.Errorf("One of 'public-key' and 'file' must be provided.")
	}

	resource.Params = &paramsUpload{
		opts: opts,
	}

	return nil
}

func (command *commandUpload) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsUpload).opts
	keypair, err := keypairs.Create(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(keypair)
}
