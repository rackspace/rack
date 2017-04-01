package usercommands

import (
	"strings"

	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/db/v1/users"
)

func singleUser(user *users.User) map[string]interface{} {
	m := structs.Map(user)

	var dbs []string
	for _, db := range user.Databases {
		dbs = append(dbs, db.Name)
	}
	m["Databases"] = strings.Join(dbs, ", ")

	return m
}
