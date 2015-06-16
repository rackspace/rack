package snapshotcommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osSnapshots "github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/rackspace/blockstorage/v1/snapshots"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists snapshots",
	Action:      commandList,
}

var keysList = []string{"ID", "Name", "Description", "Metadata", "Size", "Status", "Volume ID"}

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

	allPages, err := snapshots.List(client).AllPages()
	if err != nil {
		fmt.Printf("Error listing snapshots: %s\n", err)
		os.Exit(1)
	}
	o, err := osSnapshots.ExtractSnapshots(allPages)
	if err != nil {
		fmt.Printf("Error listing snapshots: %s\n", err)
		os.Exit(1)
	}

	f := func() interface{} {
		m := make([]map[string]interface{}, len(o))
		for j, snapshot := range o {
			m[j] = snapshotSingle(&snapshot)
		}
		return m
	}

	outputParams.F = &f
	output.Print(outputParams)

}

func snapshotSingle(rawSnapshot interface{}) map[string]interface{} {
	snapshot, ok := rawSnapshot.(osSnapshots.Snapshot)
	if !ok {
		return nil
	}

	m := structs.Map(rawSnapshot)
	m["Volume ID"] = snapshot.VolumeID

	return m

}
