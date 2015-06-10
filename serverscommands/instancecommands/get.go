package instancecommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	"github.com/olekukonko/tablewriter"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get %s [optional flags]", util.Name, commandPrefix, idOrNameUsage),
	Description: "Retrieves an existing server",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet))
	},
}

func flagsGet() []cli.Flag {
	return idAndNameFlags
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	serverID := idOrName(c, client)
	o, err := servers.Get(client, serverID).Extract()
	if err != nil {
		fmt.Printf("Error retrieving server (%s): %s\n", serverID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}

func tableGet(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{"property", "value"})
	keys := []string{"ID", "Name", "Status", "Created", "Updated", "Image", "Flavor", "Public IPv4", "Public IPv6", "Private IPv4", "KeyName"}

	m["Public IPv4"] = m["AccessIPv4"]
	m["Public IPv6"] = m["AccessIPv6"]

	// Private IPv4 case
	// Nested m["Addresses"]["private"][0]
	addrs, ok := m["Addresses"].(map[string]interface{})
	if ok {
		ips, ok := addrs["private"].([]interface{})
		if ok || len(ips) > 0 {
			priv, ok := ips[0].(map[string]interface{})
			if ok {
				m["Private IPv4"] = priv["addr"]
			}
		}
	}
	if !ok { // if any were not ok, set the field to blank
		m["Private IPv4"] = ""
	}

	flavor, ok := m["Flavor"].(map[string]interface{})
	if ok {
		m["Flavor"] = flavor["id"]
	} else {
		m["Flavor"] = ""
	}
	image, ok := m["Image"].(map[string]interface{})
	if ok {
		m["Image"] = image["id"]
	} else {
		m["Image"] = ""
	}

	for _, key := range keys {
		t.Append([]string{key, fmt.Sprint(m[key])})
	}
	t.Render()
}
