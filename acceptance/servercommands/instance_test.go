// +build acceptance

package servercommands

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"

	osServers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestServerCommands(t *testing.T) {
	fmt.Println("Creating server...")
	createdServer := createServer(t)

	defer deleteServer(t, createdServer.ID)

	fmt.Println("Updating server...")
	updateServer(t, createdServer.ID)

	fmt.Println("Retrieving server...")
	getServer(t, createdServer.ID)

	fmt.Println("Rebooting server...")
	rebootServer(t, createdServer.ID)

	fmt.Println("Rebuilding server...")
	rebuildServer(t, createdServer.ID, createdServer.AdminPass)
}

func createServer(t *testing.T) *osServers.Server {
	output, err := exec.Command("rack", "servers", "instance", "create", "--output", "json", "--image-id",
		"09de0a66-3156-48b4-90a5-1cf25a905207", "--flavor-id", "3", "--name", "rackAcceptanceTestCreated", "--wait-for-completion").Output()
	th.AssertNoErr(t, err)

	var server osServers.Server
	err = json.Unmarshal(output, &server)
	th.AssertNoErr(t, err)

	return &server
}

func getServer(t *testing.T, serverID string) {
	_, err := exec.Command("rack", "servers", "instance", "get", "--output", "json", "--id", serverID).Output()
	th.AssertNoErr(t, err)
}

func updateServer(t *testing.T, serverID string) {
	output, err := exec.Command("rack", "servers", "instance", "update", "--output", "json", "--id", serverID,
		"--rename", "rackAcceptanceTestUpdated").Output()
	th.AssertNoErr(t, err)
	var server osServers.Server
	err = json.Unmarshal(output, &server)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "rackAcceptanceTestUpdated", server.Name)
}

func deleteServer(t *testing.T, serverID string) {
	fmt.Println("Deleting server...")
	output, err := exec.Command("rack", "servers", "instance", "delete", "--output", "json", "--id", serverID, "--wait-for-completion").Output()
	th.AssertNoErr(t, err)

	type result struct {
		result string
	}
	var res result
	err = json.Unmarshal(output, &res)
	th.AssertNoErr(t, err)
}

func rebootServer(t *testing.T, serverID string) {
	output, err := exec.Command("rack", "servers", "instance", "reboot", "--output", "json", "--id", serverID, "--soft", "--wait-for-completion").Output()
	th.AssertNoErr(t, err)
	type result struct {
		result string
	}
	var res result
	err = json.Unmarshal(output, &res)
	th.AssertNoErr(t, err)
}

func rebuildServer(t *testing.T, serverID, adminPass string) {
	_, err := exec.Command("rack", "servers", "instance", "rebuild", "--output", "json", "--id", serverID, "--admin-pass", adminPass,
		"--image-id", "4315b2dc-23fc-4d81-9e73-aa620357e1d8", "--wait-for-completion").Output()
	th.AssertNoErr(t, err)
}
