package keypaircommands

import (
	"fmt"
	"io/ioutil"
	"os"

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
	Usage:       fmt.Sprintf("%s %s create <keypairName> [flags]", util.Name, commandPrefix),
	Description: "Creates a keypair",
	Action:      commandCreate,
	Flags:       util.CommandFlags(flagsCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "publicKey",
			Usage: `[optional] The public ssh key to associate with the user's account.
	It may be the actual key or the file containing the key. If empty,
	the key will be created for you and returned in the output.`,
		},
	}
}

func commandCreate(c *cli.Context) {
	util.CheckArgNum(c, 1)
	keypairName := c.Args()[0]
	client := auth.NewClient("compute")
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
	if err != nil {
		fmt.Printf("Error creating keypair [%s]: %s\n", keypairName, err)
		os.Exit(1)
	}
	output.Print(c, o, tableCreate)
}

func tableCreate(c *cli.Context, i interface{}) {
	keys := []string{"Name", "Fingerprint", "PublicKey", "PrivateKey"}
	f := func() map[string]interface{} {
		return singleKeypair(i)
	}
	output.MetadataTable(c, &f, keys)
}

func singleKeypair(i interface{}) map[string]interface{} {
	return structs.Map(i)
}
