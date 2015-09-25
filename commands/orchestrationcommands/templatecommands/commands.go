package templatecommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration template"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for a `orchestration template` request.
func Get() []cli.Command {
	return []cli.Command{
		get,
		validate,
	}
}
