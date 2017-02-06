package objectcommands

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/containers"
)

func CheckContainerExists(sc *gophercloud.ServiceClient, containerName string) error {
	containerRaw := containers.Get(sc, containerName)
	return containerRaw.Err
}
