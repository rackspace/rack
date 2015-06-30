package setup

import (
	"runtime"

	"github.com/codegangsta/cli"
)

func Init(c *cli.Context) {
	switch runtime.GOOS {
	case "linux":
		break
	case "darwin":
		break
	default:
		break
	}
}
