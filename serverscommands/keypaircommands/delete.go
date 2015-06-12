package keypaircommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       fmt.Sprintf("%s %s delete <keypairName> [flags]", util.Name, commandPrefix),
	Description: "Deletes a keypair",
	Action:      commandDelete,
	Flags:       util.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{}
}

var keysDelete = []string{}

func commandDelete(c *cli.Context) {
	util.CheckArgNum(c, 1)
	keypairName := c.Args()[0]
	client := auth.NewClient("compute")
	err := keypairs.Delete(client, keypairName).ExtractErr()
	if err != nil {
		fmt.Printf("Error deleting keypair [%s]: %s\n", keypairName, err)
		os.Exit(1)
	}
}
