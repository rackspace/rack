package keypaircommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists keypairs",
	Action:      commandList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{}
}

var keysList = []string{"Name", "Fingerprint"}

func commandList(c *cli.Context) {
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysList,
	}
	err = util.CheckArgNum(c, 0)
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

	allPages, err := keypairs.List(client).AllPages()
	outputParams.ServiceClient = client
	if err != nil {
		fmt.Printf("Error listing keypairs: %s\n", err)
		os.Exit(1)
	}
	o, err := osKeypairs.ExtractKeyPairs(allPages)
	if err != nil {
		fmt.Printf("Error listing keypairs: %s\n", err)
		os.Exit(1)
	}

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, kp := range o {
			m[j] = structs.Map(kp)
		}
		return m
	}
	outputParams.F = &f
	output.Print(outputParams)
}
