package stacktemplatecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration stack-template"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for an `orchestration template` request.
func Get() []cli.Command {
	return []cli.Command{
		get,
		validate,
	}
}
