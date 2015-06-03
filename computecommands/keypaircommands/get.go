package keypaircommands

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/gophercloud/rackspace/compute/v2/keypairs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
)

var get = cli.Command{
	Name:        "get",
	Usage:       fmt.Sprintf("%s %s get <keypairName> [flags]", util.Name, commandPrefix),
	Description: "Retreives a keypair",
	Action:      commandGet,
	Flags:       flagsGet(),
}

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

func commandGet(c *cli.Context) {
	util.CheckArgNum(c, 1)
	flavorID := c.Args()[0]
	client := auth.NewClient("compute")
	o, err := keypairs.Get(client, flavorID).Extract()
	if err != nil {
		fmt.Printf("Error retreiving image [%s]: %s\n", flavorID, err)
		os.Exit(1)
	}
	output.Print(c, o, tableGet)
}

func tableGet(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	colWidth := 100
	t.SetColWidth(colWidth)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{"property", "value"})
	keys := []string{"Name", "Fingerprint"}
	for _, key := range keys {
		t.Append([]string{key, fmt.Sprint(m[key])})
	}
	t.Render()
	keyWidth := tablewriter.DisplayWidth(m["PublicKey"].(string))
	numPieces := keyWidth / colWidth
	remainder := keyWidth % colWidth
	pieces := make([]string, numPieces)
	j := 0
	for j = 0; j < numPieces; {
		pieces = append(pieces, m["PublicKey"].(string)[colWidth*j:colWidth*j+colWidth])
		j++
	}
	pieces = append(pieces, m["PublicKey"].(string)[colWidth*j:colWidth*j+remainder])
	pk := strings.TrimLeft(strings.Join(pieces, "\n"), "\n")
	fmt.Printf("PublicKey: %s\n", pk)
}
