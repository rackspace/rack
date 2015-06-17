package instancecommands

import (
	"fmt"
	"io/ioutil"
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
	Usage:       util.Usage(commandPrefix, "create", "--name <serverName>"),
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
			Name:  "image-id",
			Usage: "[optional; required if imageName and bootFromVolume flags are not provided] The image ID from which to create the server.",
		},
		cli.StringFlag{
			Name:  "image-name",
			Usage: "[optional; required if imageRef and bootFromVolume flags are not provided] The name of the image from which to create the server.",
		},
		cli.StringFlag{
			Name:  "flavor-id",
			Usage: "[optional; required if flavorName is not provided] The flavor ID that the server should have.",
		},
		cli.StringFlag{
			Name:  "flavor-name",
			Usage: "[optional; required if flavorRef is not provided] The name of the flavor that the server should have.",
		},
		cli.StringFlag{
			Name:  "security-groups",
			Usage: "[optional] A comma-separated string of names of the security groups to which this server should belong.",
		},
		cli.StringFlag{
			Name:  "user-data",
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
			Name:  "admin-pass",
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
	var err error
	outputParams := &output.Params{
		Context: c,
		Keys:    keysCreate,
	}
	err = util.CheckArgNum(c, 0)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	err = util.CheckFlagsSet(c, []string{"name"})
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	serverName := c.String("name")

	opts := &servers.CreateOpts{
		Name:           serverName,
		ImageRef:       c.String("image-id"),
		ImageName:      c.String("image-name"),
		FlavorRef:      c.String("flavor-id"),
		FlavorName:     c.String("flavor-name"),
		SecurityGroups: strings.Split(c.String("security-groups"), ","),
		AdminPass:      c.String("admin-pass"),
		KeyPair:        c.String("keypair"),
	}

	if c.IsSet("user-data") {
		s := c.String("user-data")
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
		opts.Metadata, err = util.CheckKVFlag(c, "metadata")
		if err != nil {
			outputParams.Err = err
			output.Print(outputParams)
			return
		}
	}

	outputParams.ServiceClientType = serviceClientType
	client, err := auth.NewClient(c, outputParams.ServiceClientType)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	o, err := servers.Create(client, opts).Extract()
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error creating server (%s): %s\n", serverName, err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		return serverSingle(o)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
