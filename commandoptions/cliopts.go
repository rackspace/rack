package commandoptions

import "github.com/codegangsta/cli"

type Cred struct {
	Value string
	From  string
}

func CLIopts(c *cli.Context, have map[string]Cred, need map[string]string) {
	for opt := range need {
		if c.IsSet(opt) {
			have[opt] = Cred{Value: c.String(opt), From: "command-line"}
			delete(need, opt)
		}
	}
}
