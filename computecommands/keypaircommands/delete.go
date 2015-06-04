package keypaircommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/util"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       fmt.Sprintf("%s %s delete <keypairName> [flags]", util.Name, commandPrefix),
	Description: "Deletes a keypair",
	Action:      commandDelete,
	Flags:       util.CommandFlags(flagsDelete),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{}
}

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
