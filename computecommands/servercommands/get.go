package servercommands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get <serverID> [flags]", util.Name, commandPrefix),
	Description: "Retrieves an existing server",
	Action:      commandGet,
	Flags:       flagsGet(),
}

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 1)
	serverID := c.Args()[0]
	client := auth.NewClient("compute")
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
	for _, key := range keys {
		tmp := ""
		switch key {
		case "Public IPv4":
			tmp = fmt.Sprint(m["AccessIPv4"])
		case "Public IPv6":
			tmp = fmt.Sprint(m["AccessIPv6"])
		case "Private IPv4":
			i, ok := m["Addresses"].(map[string]interface{})
			if !ok {
				tmp = ""
				break
			}
			j, ok := i["private"].([]interface{})
			if !ok || len(j) == 0 {
				tmp = ""
				break
			}
			i, ok = j[0].(map[string]interface{})
			if !ok {
				tmp = ""
				break
			}
			tmp = fmt.Sprint(i["addr"])
		case "Image":
			i, ok := m["Image"].(map[string]interface{})
			if !ok {
				tmp = ""
				break
			}
			tmp = fmt.Sprint(i["id"])
		case "Flavor":
			i, ok := m["Flavor"].(map[string]interface{})
			if !ok {
				tmp = ""
				break
			}
			tmp = fmt.Sprint(i["id"])
		default:
			tmp = fmt.Sprint(m[key])
		}
		if tmp == "<nil>" {
			tmp = ""
		}
		t.Append([]string{key, tmp})
	}
	t.Render()
}
