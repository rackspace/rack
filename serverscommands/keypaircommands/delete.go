package keypaircommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "--name <keypairName>"),
	Description: "Deletes a keypair",
	Action:      commandDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the keypair",
		},
	}
}

var keysDelete = []string{}

func commandDelete(c *cli.Context) {
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
	err = keypairs.Delete(client, keypairName).ExtractErr()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error deleting keypair [%s]: %s\n", keypairName, err)
		output.Print(outputParams)
		return
	}
	f := func() interface{} {
		return fmt.Sprintf("Successfully deleted keypair [%s]", keypairName)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
