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
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
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
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	allPages, err := keypairs.List(client).AllPages()
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
	output.Print(c, &f, keysList)
}
