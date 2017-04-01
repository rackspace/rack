package instancecommands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/instances"
)

func singleInstance(rawInstance interface{}) map[string]interface{} {
	instance, ok := rawInstance.(*instances.Instance)
	if !ok {
		return nil
	}

	m := structs.Map(rawInstance)
	m["Datastore"] = fmt.Sprintf("%s %s", instance.Datastore.Type, instance.Datastore.Version)
	m["Flavor"] = instance.Flavor.ID
	m["Size"] = fmt.Sprintf("%s GB", strconv.Itoa(instance.Volume.Size))
	m["Created"] = instance.Created.Format(time.RFC1123)
	m["Updated"] = instance.Updated.Format(time.RFC1123)

	return m
}

func singleUser(rawUser interface{}) map[string]interface{} {
	user, ok := rawUser.(*users.User)
	if !ok {
		return nil
	}

	return structs.Map(user)
}
