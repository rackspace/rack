package backupcommands

import (
	"fmt"
	"strconv"

	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/instances"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/backups"
)

func singleBackup(backup *backups.Backup) map[string]interface{} {
	return structs.Map(backup)
}

func singleInstance(rawInstance interface{}) map[string]interface{} {
	instance, ok := rawInstance.(*instances.Instance)
	if !ok {
		return nil
	}

	m := structs.Map(rawInstance)
	m["Datastore"] = fmt.Sprintf("%s %s", instance.Datastore.Type, instance.Datastore.Version)
	m["Flavor"] = instance.Flavor.ID
	m["Size"] = fmt.Sprintf("%s GB", strconv.Itoa(instance.Volume.Size))

	return m
}
