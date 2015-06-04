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
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [flags]", util.Name, commandPrefix),
	Description: "Lists existing servers",
	Action:      commandList,
	Flags:       flagsList(),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(flagsList)
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list servers with this name.",
		},
		cli.StringFlag{
			Name:  "changesSince",
			Usage: "Only list servers that have been changed since this time/date stamp.",
		},
		cli.StringFlag{
			Name:  "image",
			Usage: "Only list servers that have this image ID.",
		},
		cli.StringFlag{
			Name:  "flavor",
			Usage: "Only list servers that have this flavor ID.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list servers that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "Start listing servers at this server ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "Only return this many servers at most.",
		},
	}
}

func commandList(c *cli.Context) {
	util.CheckArgNum(c, 0)
	client := auth.NewClient("compute")
	opts := osServers.ListOpts{
		ChangesSince: c.String("changesSince"),
		Image:        c.String("image"),
		Flavor:       c.String("flavor"),
		Name:         c.String("name"),
		Status:       c.String("status"),
		Marker:       c.String("marker"),
		Limit:        c.Int("limit"),
	}
	allPages, err := servers.List(client, opts).AllPages()
	if err != nil {
		fmt.Printf("Error listing servers: %s\n", err)
		os.Exit(1)
	}
	o, err := servers.ExtractServers(allPages)
	if err != nil {
		fmt.Printf("Error listing servers: %s\n", err)
		os.Exit(1)
	}
	output.Print(c, o, tableList)
}

func tableList(c *cli.Context, i interface{}) {
	servers, ok := i.([]osServers.Server)
	if !ok {
		fmt.Fprintf(c.App.Writer, "Could not type assert interface\n%+v\nto []osServers.Server\n", i)
		os.Exit(1)
	}
	t := tablewriter.NewWriter(c.App.Writer)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	keys := []string{"ID", "Name", "Status", "Public IPv4", "Private IPv4", "Image", "Flavor"}
	t.SetHeader(keys)
	for _, server := range servers {
		m := structs.Map(server)
		f := []string{}
		for _, key := range keys {
			tmp := ""
			switch key {
			case "Public IPv4":
				tmp = fmt.Sprint(m["AccessIPv4"])
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
			f = append(f, tmp)
		}
		t.Append(f)
	}
	t.Render()
}
