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
	Usage:       fmt.Sprintf("%s %s delete [--name <keypairName>] [flags]", util.Name, commandPrefix),
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
	util.CheckArgNum(c, 0)
	if !c.IsSet("name") {
		util.PrintError(c, util.ErrMissingFlag{
			Msg: "--name is required.",
		})
	}
	kpName := c.String("name")
	client := auth.NewClient("compute")
	err := keypairs.Delete(client, kpName).ExtractErr()
	if err != nil {
		fmt.Printf("Error deleting keypair [%s]: %s\n", kpName, err)
		os.Exit(1)
	}
}
