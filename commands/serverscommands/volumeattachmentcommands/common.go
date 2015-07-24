package volumeattachmentcommands

import (
	"fmt"

	"github.com/jrperritt/rack/handler"
	osServers "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/jrperritt/rack/output"
)

func serverIDorName(ctx *handler.Context) (string, error) {
	if ctx.CLIContext.IsSet("server-id") {
		if ctx.CLIContext.IsSet("server-name") {
			return "", fmt.Errorf("Only one of either --server-id or --server-name may be provided.")
		}
		return ctx.CLIContext.String("server-id"), nil
	} else if ctx.CLIContext.IsSet("server-name") {
		name := ctx.CLIContext.String("server-name")
		id, err := osServers.IDFromName(ctx.ServiceClient, name)
		if err != nil {
			return "", fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		return id, nil
	} else {
		return "", output.ErrMissingFlag{Msg: "One of either --server-id or --server-name must be provided."}
	}
}
