package instancecommands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var list = cli.Command{
	Name:        "list",
	Usage:       fmt.Sprintf("%s %s list [optional flags]", util.Name, commandPrefix),
	Description: "Lists existing servers",
	Action:      commandList,
	Flags:       util.CommandFlags(flagsList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList))
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

	keys := []string{"ID", "Name", "Status", "Public IPv4", "Private IPv4", "Image", "Flavor"}

	w := new(tabwriter.Writer)
	w.Init(c.App.Writer, 0, 8, 0, '\t', 0)

	// Write the header
	fmt.Fprintln(w, strings.Join(keys, "\t"))
	for _, server := range servers {
		m := structs.Map(server)

		// Extract the Image ID - if !ok, image is ""
		image, _ := GetNestedID(m["Image"])

		// Extract the Flavor ID - if !ok, flavor is ""
		flavor, _ := GetNestedID(m["Flavor"])

		// Extract the very first private address
		// TODO: How do we handle multiples here?
		privAddr := ""
		a, ok := m["Addresses"].(map[string]interface{})
		if ok {
			a, ok := a["private"].([]interface{})
			if ok && len(a) > 0 {
				first, ok := a[0].(map[string]interface{})
				if ok {
					privAddr = fmt.Sprint(first["addr"])
				}

			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", m["ID"], m["Name"], m["Status"], m["AccessIPv4"], privAddr, image, flavor)

	}
	w.Flush()
}

// GetNestedID extracts the nested id from, e.g. the flavor from a server
func GetNestedID(i interface{}) (string, bool) {
	f, ok := i.(map[string]interface{})
	if !ok {
		return "", ok
	}
	return fmt.Sprint(f["id"]), ok
}
