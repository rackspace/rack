package stackresourcecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration resource"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for a `orch resource` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		getMetadata,
		getSchema,
		getTemplate,
		listTypes,
	}
}
