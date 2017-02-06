package objectcommands

import "github.com/codegangsta/cli"

var commandPrefix = "files object"
var serviceClientType = "object-store"

// Get returns all the commands allowed for a `files object` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		upload,
		uploadDir,
		download,
		get,
		remove,
		setMetadata,
		updateMetadata,
		getMetadata,
		deleteMetadata,
	}
}
