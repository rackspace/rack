package buildinfocommands

import "github.com/rackspace/rack/internal/github.com/codegangsta/cli"

var commandPrefix = "orchestration buildinfo"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for a `orch event` request.
func Get() []cli.Command {
	return []cli.Command{
		get,
	}
}
