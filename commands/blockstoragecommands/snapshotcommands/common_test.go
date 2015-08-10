package snapshotcommands

import (
	"testing"

	osSnapshots "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
)

func TestSnapshotSingle(t *testing.T) {
	snapshot := &osSnapshots.Snapshot{
		ID:          "47e3dd45-f725-4183-81be-cfe97f2c2a08",
		Name:        "test-rack-snapshot",
		Size:        500,
		Status:      "available",
		VolumeID:    "7ba4e908-8a01-462d-b4a4-2f86ba8c2432",
		Attachments: []string{"123456789", "987654321"},
		Bootable:    "true",
	}

	actual := snapshotSingle(snapshot)

	th.AssertEquals(t, "123456789,987654321", actual["Attachments"])
}
