package keypaircommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get <keypairName> [flags]", util.Name, commandPrefix),
	Description: "Retreives a keypair",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 1)
	flavorID := c.Args()[0]
	client := auth.NewClient("compute")
	o, err := keypairs.Get(client, flavorID).Extract()
	if err != nil {
		fmt.Printf("Error retreiving image [%s]: %s\n", flavorID, err)
		os.Exit(1)
	}

	// Assume they want the key directly
	fmt.Fprintf(c.App.Writer, "%s", o.PublicKey)
}
