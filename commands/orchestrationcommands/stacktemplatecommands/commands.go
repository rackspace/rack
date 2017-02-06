package stacktemplatecommands

import "github.com/codegangsta/cli"

var commandPrefix = "orchestration template"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for an `orchestration template` request.
func Get() []cli.Command {
	return []cli.Command{
		validate,
	}
}
