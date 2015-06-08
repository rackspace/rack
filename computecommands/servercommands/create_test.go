package servercommands

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/jrperritt/rack/util"
)

var baseCmd = fmt.Sprintf("rack %s create", commandPrefix)

func TestServerCreateNoName(t *testing.T) {
	cmd := baseCmd
	cmdSlice := strings.Split(cmd, " ")
	out, _ := exec.Command(cmdSlice[0], cmdSlice[1:]...).Output()
	if !strings.HasPrefix(string(out), util.ErrMissingFlagPrefix) {
		t.Log("Expected error due to empty --name but didn't get one.")
		t.Fail()
	}
}

func TestServerCreateInvalidMetadata(t *testing.T) {
	cmd := baseCmd + " --name foo --metadata foo=bar;a=b"
	cmdSlice := strings.Split(cmd, " ")
	out, _ := exec.Command(cmdSlice[0], cmdSlice[1:]...).Output()
	if !strings.HasPrefix(string(out), util.ErrFlagFormattingPrefix) {
		t.Log("Expected error due to invalid --metadata format but didn't get one.")
		t.Fail()
	}
}
