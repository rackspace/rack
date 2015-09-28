package stackresourcecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration stack-resource"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for a `orchestration resource` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		getSchema,
		getTemplate,
		listTypes,
	}
}
