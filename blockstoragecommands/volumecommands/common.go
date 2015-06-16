package volumecommands

import (
	"github.com/fatih/structs"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func volumeSingle(rawVolume interface{}) map[string]interface{} {
	volume, ok := rawVolume.(*volumes.Volume)
	if !ok {
		return nil
	}

	m := structs.Map(rawVolume)
	m["Volume Type"] = volume.VolumeType
	m["Snapshot ID"] = volume.SnapshotID
	m["Created"] = volume.CreatedAt

	return m

}
