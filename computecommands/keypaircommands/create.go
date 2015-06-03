package keypaircommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rackcli/auth"
	"github.com/jrperritt/rackcli/output"
	"github.com/jrperritt/rackcli/util"
	"github.com/olekukonko/tablewriter"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var create = cli.Command{
	Name:        "create",
	Usage:       fmt.Sprintf("%s %s create <keypairName> [flags]", util.Name, commandPrefix),
	Description: "Creates a keypair",
	Action:      commandCreate,
	Flags:       flagsCreate(),
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "publicKey",
			Usage: `[optional] The public ssh key to associate with the user's account.
	It may be the actual key or the file containing the key. If empty,
	the key will be created for you and returned in the output.`,
		},
	}
}

func commandCreate(c *cli.Context) {
	util.CheckArgNum(c, 1)
	keypairName := c.Args()[0]
	client := auth.NewClient("compute")
	opts := osKeypairs.CreateOpts{
		Name: keypairName,
	}

	if c.IsSet("publicKey") {
		s := c.String("publicKey")
		pk, err := ioutil.ReadFile(s)
		if err != nil {
			opts.PublicKey = string(pk)
		} else {
			opts.PublicKey = s
		}
	}

	o, err := keypairs.Create(client, opts).Extract()
	if err != nil {
		fmt.Printf("Error creating keypair [%s]: %s\n", keypairName, err)
		os.Exit(1)
	}
	output.Print(c, o, tableCreate)
}

func tableCreate(c *cli.Context, i interface{}) {
	m := structs.Map(i)
	t := tablewriter.NewWriter(c.App.Writer)
	colWidth := 100
	t.SetColWidth(colWidth)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader([]string{"property", "value"})
	keys := []string{"Name", "Fingerprint", "PublicKey", "PrivateKey"}
	for _, key := range keys {
		switch key {
		case "PublicKey", "PrivateKey":
			keyWidth := tablewriter.DisplayWidth(m[key].(string))
			numPieces := keyWidth / colWidth
			remainder := keyWidth % colWidth
			pieces := make([]string, numPieces)
			j := 0
			for j = 0; j < numPieces; {
				pieces = append(pieces, m[key].(string)[colWidth*j:colWidth*j+colWidth])
				j++
			}
			pieces = append(pieces, m[key].(string)[colWidth*j:colWidth*j+remainder])
			t.Append([]string{key, strings.Join(pieces, "\n")})
		default:
			t.Append([]string{key, fmt.Sprint(m[key])})
		}
	}
	t.Render()
}
