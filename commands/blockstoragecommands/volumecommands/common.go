package volumecommands

import (
	"github.com/jrperritt/rack/internal/github.com/fatih/structs"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func volumeSingle(volume *volumes.Volume) map[string]interface{} {
	m := structs.Map(volume)
	for k, v := range m {
		if v == nil {
			m[k] = ""
		}
	}
	return m
}
