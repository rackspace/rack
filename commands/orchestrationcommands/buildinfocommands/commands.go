package buildinfocommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration build-info"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for an `orchestration build-info` request.
func Get() []cli.Command {
	return []cli.Command{
		get,
	}
}
