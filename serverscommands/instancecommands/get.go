package instancecommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jrperritt/rack/auth"
	"github.com/jrperritt/rack/output"
	"github.com/jrperritt/rack/util"
	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/servers"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", util.IDOrNameUsage("instance")),
	Description: "Retrieves an existing server",
	Action:      commandGet,
	Flags:       util.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return util.IDAndNameFlags
}

var keysGet = []string{"ID", "Name", "Status", "Created", "Updated", "Image", "Flavor", "Public IPv4", "Public IPv6", "Private IPv4", "KeyName"}

func commandGet(c *cli.Context) {
	var err error
	// initialize parameters for processing the output
	outputParams := &output.Params{
		Context: c,
		Keys:    keysGet,
	}
	// make sure no arguments are passed
	err = util.CheckArgNum(c, 0)
	// if arguments were passed, return early with error message
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	// set the type of the service client. this is used in `Print` to form the
	// cache key
	outputParams.ServiceClientType = serviceClientType
	client, err := auth.NewClient(c, outputParams.ServiceClientType)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}

	// check if the user provided an ID or a name for the instance
	serverID, err := util.IDOrName(c, client, osServers.IDFromName)
	if err != nil {
		outputParams.Err = err
		output.Print(outputParams)
		return
	}
	o, err := servers.Get(client, serverID).Extract()
	// set the service client on outputParams. this value will contain the new
	// (and valid) token if the client had to re-authenticate
	outputParams.ServiceClient = client
	if err != nil {
		outputParams.Err = fmt.Errorf("Error retrieving server (%s): %s\n", serverID, err)
		output.Print(outputParams)
		return
	}

	f := func() interface{} {
		return serverSingle(o)
	}
	outputParams.F = &f
	output.Print(outputParams)
}
