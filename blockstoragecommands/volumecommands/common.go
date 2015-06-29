package volumecommands

import (
	"github.com/fatih/structs"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func volumeSingle(volume *volumes.Volume) map[string]interface{} {
	m := structs.Map(volume)
	return m
}
