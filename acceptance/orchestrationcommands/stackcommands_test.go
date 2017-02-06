// +build acceptance

package orchestrationcommands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	osStacks "github.com/rackspace/gophercloud/openstack/orchestration/v1/stacks"
	th "github.com/rackspace/gophercloud/testhelper"
)

const createTemplateURL = "https://raw.githubusercontent.com/rackerlabs/heat-ci/master/dev/smoke.yaml"
const updateTemplateURL = createTemplateURL

func TestStackCommands(t *testing.T) {
	fmt.Println("Previewing stack...")
	previewStack(t)

	fmt.Println("Creating stack...")
	createdStack := createStack(t)

	fmt.Printf("Retrieving stack %s...\n", createdStack.Result.ID)
	getStack(t, createdStack.Result.ID)

	fmt.Printf("Updating stack %s...\n", createdStack.Result.ID)
	updateStack(t, createdStack.Result.ID)

	fmt.Printf("Abandoning stack...%s\n", createdStack.Result.ID)
	adoptData := abandonStack(t, createdStack.Result.ID)

	fmt.Println("Adopting stack...")
	adoptedStack := adoptStack(t, adoptData)

	fmt.Printf("Deleting stack...%s\n", adoptedStack.Result.ID)
	deleteStack(t, adoptedStack.Result.ID)

}

func createStack(t *testing.T) *stackCreatedResponse {
	output, err := exec.Command("rack", "orchestration", "stack", "create", "--output", "json", "--template-url",
		createTemplateURL, "--name", "rackAcceptanceTestStackCreated").Output()
	th.AssertNoErr(t, err)
	var stack stackCreatedResponse
	err = json.Unmarshal(output, &stack)
	th.AssertNoErr(t, err)

	return &stack
}

func previewStack(t *testing.T) {
	output, err := exec.Command("rack", "orchestration", "stack", "preview", "--output", "json", "--template-url",
		createTemplateURL, "--name", "rackAcceptanceTestStackCreated").Output()
	th.AssertNoErr(t, err)
	var stack stackPreviewedResponse
	err = json.Unmarshal(output, &stack)
	th.AssertNoErr(t, err)
}

func getStack(t *testing.T, stackID string) {
	_, err := exec.Command("rack", "orchestration", "stack", "get", "--output", "json", "--id", stackID).Output()
	th.AssertNoErr(t, err)
}

func updateStack(t *testing.T, stackID string) {
	output, err := exec.Command("rack", "orchestration", "stack", "update", "--output", "json", "--id", stackID,
		"--template-url", updateTemplateURL).Output()
	th.AssertNoErr(t, err)
	var stack stackUpdatedResponse
	err = json.Unmarshal(output, &stack)
	th.AssertNoErr(t, err)
}

func deleteStack(t *testing.T, stackID string) {
	_, err := exec.Command("rack", "orchestration", "stack", "delete", "--id", stackID).Output()
	th.AssertNoErr(t, err)
}

func abandonStack(t *testing.T, stackID string) []byte {
	output, err := exec.Command("rack", "orchestration", "stack", "abandon", "--output", "json", "--id", stackID).Output()
	th.AssertNoErr(t, err)
	return output
}

func adoptStack(t *testing.T, adoptData []byte) *stackCreatedResponse {
	// write adoptData to a temporary file
	adoptFile, err := ioutil.TempFile(os.TempDir(), "adoptFile")
	th.AssertNoErr(t, err)
	defer os.Remove(adoptFile.Name())
	adoptFile.Write(adoptData)

	output, err := exec.Command("rack", "orchestration", "stack", "adopt", "--adopt-file", adoptFile.Name(), "--name", "rackAcceptanceTestStackCreated", "--output", "json").Output()
	th.AssertNoErr(t, err)
	var stack stackCreatedResponse
	err = json.Unmarshal(output, &stack)
	th.AssertNoErr(t, err)
	return &stack
}

type stackCreatedResponse struct {
	Result osStacks.CreatedStack `json:"result"`
}

type stackUpdatedResponse struct {
	Result osStacks.RetrievedStack `json:"result"`
}

type stackPreviewedResponse struct {
	Result osStacks.RetrievedStack `json:"result"`
}
