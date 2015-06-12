package keypaircommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s [globals] %s get [--name <keypairName>] [flags]", util.Name, commandPrefix),
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
	util.CheckArgNum(c, 0)
	if !c.IsSet("name") {
		util.PrintError(c, util.ErrMissingFlag{
			Msg: "--name is required.",
		})
	}
	kpName := c.String("name")
	client := auth.NewClient("compute")
	o, err := keypairs.Get(client, kpName).Extract()
	if err != nil {
		fmt.Printf("Error retreiving keypair [%s]: %s\n", kpName, err)
		os.Exit(1)
	}

	f := func() interface{} {
		m := structs.Map(o)
		if c.GlobalIsSet("json") {
			return m
		}
		// Assume they want the key directly
		return m["PublicKey"]
	}
	output.Print(c, &f, keysGet)
}
