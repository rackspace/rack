package volumecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osVolumes "github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/rackspace/gophercloud/rackspace/blockstorage/v1/volumes"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists volumes",
	Action:      commandList,
}

var keysList = []string{"ID", "Name", "Description", "Size", "Volume Type", "Snapshot ID", "Attachments", "Created"}

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

	allPages, err := volumes.List(client).AllPages()
	if err != nil {
		fmt.Printf("Error listing volumes: %s\n", err)
		os.Exit(1)
	}
	o, err := osVolumes.ExtractVolumes(allPages)
	if err != nil {
		fmt.Printf("Error listing volumes: %s\n", err)
		os.Exit(1)
	}

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, volume := range o {
			m[j] = volumeSingle(&volume)
		}
		return m
	}

	outputParams.F = &f
	output.Print(outputParams)
}
