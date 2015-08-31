package objectcommands

import (
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
)

func checkContainerExists(sc *gophercloud.ServiceClient, containerName string) error {
	containerRaw := containers.Get(sc, containerName)
	return containerRaw.Err
}
