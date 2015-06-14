package instancecommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var create = cli.Command{
	Name:        "create",
	Usage:       fmt.Sprintf("%s %s create [--name <serverName>] [optional flags]", util.Name, commandPrefix),
	Description: "Creates a new server",
	Action:      commandCreate,
	Flags:       util.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name that the server should have.",
		},
		cli.StringFlag{
			Name:  "imageRef",
			Usage: "[optional; required if imageName and bootFromVolume flags are not provided] The image ID from which to create the server.",
		},
		cli.StringFlag{
			Name:  "imageName",
			Usage: "[optional; required if imageRef and bootFromVolume flags are not provided] The name of the image from which to create the server.",
		},
		cli.StringFlag{
			Name:  "flavorRef",
			Usage: "[optional; required if flavorName is not provided] The flavor ID that the server should have.",
		},
		cli.StringFlag{
			Name:  "flavorName",
			Usage: "[optional; required if flavorRef is not provided] The name of the flavor that the server should have.",
		},
		cli.StringFlag{
			Name:  "securityGroups",
			Usage: "[optional] A comma-separated string of names of the security groups to which this server should belong.",
		},
		cli.StringFlag{
			Name:  "userData",
			Usage: "[optional] Configuration information or scripts to use after the server boots.",
		},
		cli.StringFlag{
			Name:  "networks",
			Usage: "[optional] A comma-separated string of IDs of the networks to attach to this server. If not provided, a public and private network will be attached.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string a key=value pairs.",
		},
		cli.StringFlag{
			Name:  "adminPass",
			Usage: "[optional] The root password for the server. If not provided, one will be randomly generated and returned in the output.",
		},
		cli.StringFlag{
			Name:  "keypair",
			Usage: "[optional] The name of the already-existing SSH KeyPair to be injected into this server.",
		},
	}
}

var keysCreate = []string{"ID", "AdminPass"}

func commandCreate(c *cli.Context) {
	util.CheckArgNum(c, 0)

	if !c.IsSet("name") {
		util.PrintError(c, util.ErrMissingFlag{
			Msg: "--name is required.",
		})
	}

	opts := &servers.CreateOpts{
		Name:           c.String("name"),
		ImageRef:       c.String("imageRef"),
		ImageName:      c.String("imageName"),
		FlavorRef:      c.String("flavorRef"),
		FlavorName:     c.String("flavorName"),
		SecurityGroups: strings.Split(c.String("securityGroups"), ","),
		AdminPass:      c.String("adminPass"),
		KeyPair:        c.String("keypair"),
	}

	if c.IsSet("userData") {
		s := c.String("userData")
		userData, err := ioutil.ReadFile(s)
		if err != nil {
			opts.UserData = userData
		} else {
			opts.UserData = []byte(s)
		}
	}

	if c.IsSet("networks") {
		netIDs := strings.Split(c.String("networks"), ",")
		networks := make([]osServers.Network, len(netIDs))
		for i, netID := range netIDs {
			networks[i] = osServers.Network{
				UUID: netID,
			}
		}
		opts.Networks = networks
	}

	if c.IsSet("metadata") {
		opts.Metadata = util.CheckKVFlag(c, "metadata")
	}

	client := auth.NewClient("compute")
	o, err := servers.Create(client, opts).Extract()
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		os.Exit(1)
	}

	f := func() interface{} {
		return serverSingle(o)
	}
	output.Print(c, &f, keysCreate)
}
