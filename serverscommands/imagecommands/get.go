package imagecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/images"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get <imageID> [flags]", util.Name, commandPrefix),
	Description: "Retreives an image",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

var keysGet = []string{"ID", "Name", "Status", "Progress", "MinDisk", "MinRAM", "Created", "Updated"}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 1)
	imageID := c.Args()[0]
	client := auth.NewClient("compute")
	o, err := images.Get(client, imageID).Extract()
	if err != nil {
		fmt.Printf("Error retreiving image [%s]: %s\n", imageID, err)
		os.Exit(1)
	}
	f := func() interface{} {
		return structs.Map(o)
	}
	output.Print(c, &f, keysGet)
}
