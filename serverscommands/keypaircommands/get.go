package keypaircommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "--name <keypairName>"),
	Description: "Retreives a keypair",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the keypair",
		},
	}
}

var keysGet = []string{"Name", "Fingerprint", "PublicKey", "UserID"}

func commandGet(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysGet,
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
	o, err := keypairs.Get(client, keypairName).Extract()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error retreiving keypair [%s]: %s\n", keypairName, err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		m := structs.Map(o)
		if c.GlobalIsSet("json") {
			return m
		}
		// Assume they want the key directly
		return m["PublicKey"]
	}
	outputParams.F = &f
	output.Print(outputParams)
}
