// +build acceptance

package filescommands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"testing"

	th "github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/testhelper"
)

func TestLargeObjectCommands(t *testing.T) {
	containerName := "rackAcceptanceTestLargeObjectContainer"
	objectName := "rackAcceptanceTestLargeObject"
	var sizeFile int64 = 250
	sizePieces := 50

	fmt.Println("Creating container...")
	createContainer(t, containerName)

	defer emptyAndDeleteContainer(t, containerName)

	fmt.Println("Creating and uploading large object...")
	uploadLargeObject(t, containerName, objectName, sizeFile, sizePieces)
}

func createContainer(t *testing.T, containerName string) {
	_, err := exec.Command("rack", "files", "container", "create", "--output", "json",
		"--name", containerName).Output()
	th.AssertNoErr(t, err)
}

func uploadLargeObject(t *testing.T, containerName, objectName string, sizeFile int64, sizePieces int) {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	sizeFile = int64(sizeFile * 1000000)

	tempFile, err := ioutil.TempFile(".", "tmp")
	th.AssertNoErr(t, err)

	defer func() {
		err := os.Remove(tempFile.Name())
		th.AssertNoErr(t, err)
	}()

	data := make([]byte, sizeFile)
	for i := range data {
		data[i] = byte(letters[rand.Intn(len(letters))])
	}

	err = ioutil.WriteFile(tempFile.Name(), data, 0777)
	th.AssertNoErr(t, err)

	_, err = exec.Command("rack", "files", "large-object", "upload", "--output", "json",
		"--container", containerName, "--name", objectName, "--file", tempFile.Name(), "--size-pieces", strconv.Itoa(sizePieces)).Output()
	th.AssertNoErr(t, err)

}

func emptyAndDeleteContainer(t *testing.T, containerName string) {
	_, err := exec.Command("rack", "files", "container", "empty", "--name", containerName).Output()
	th.AssertNoErr(t, err)

	_, err := exec.Command("rack", "files", "container", "delete", "--name", containerName).Output()
	th.AssertNoErr(t, err)
}
