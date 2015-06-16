package keypaircommands

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--name <keypairName>"),
	Description: "Creates a keypair",
	Action:      commandCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "publicKey",
			Usage: strings.Join([]string{"[optional] The public ssh key to associate with the user's account.",
				"It may be the actual key or the file containing the key. If empty,",
				"the key will be created for you and returned in the output."}, "\n\t"),
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the keypair",
		},
	}
}

var keysCreate = []string{"Name", "Fingerprint", "PublicKey", "PrivateKey"}

func commandCreate(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysCreate,
	}
	err = util.CheckArgNum(c, 0)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	err = util.CheckFlagsSet(c, []string{"name"})
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	outputParams.ServiceClientType = serviceClientType
	client, err := auth.NewClient(c, outputParams.ServiceClientType)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	keypairName := c.String("name")
	opts := osKeypairs.CreateOpts{
		Name: keypairName,
	}

	if c.IsSet("publicKey") {
		s := c.String("publicKey")
		pk, err := ioutil.ReadFile(s)
		if err != nil {
			opts.PublicKey = string(pk)
		} else {
			opts.PublicKey = s
		}
	}

	o, err := keypairs.Create(client, opts).Extract()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error creating keypair [%s]: %s\n", keypairName, err)
		output.Print(outputParams)
		return
	}
	f := func() interface{} {
		return structs.Map(o)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
