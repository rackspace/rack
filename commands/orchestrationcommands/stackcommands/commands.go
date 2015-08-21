package stackcommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orch stacks"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for a `orch build-info` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
	}
}
