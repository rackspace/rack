package commandoptions

import "github.com/jrperritt/rack/internal/github.com/codegangsta/cli"

type Cred struct {
	Value string
	From  string
}

func CLIopts(c *cli.Context, have map[string]Cred, need map[string]string) {
	for opt := range need {
		if c.GlobalIsSet(opt) {
			have[opt] = Cred{Value: c.GlobalString(opt), From: "command-line"}
			delete(need, opt)
		} else if c.IsSet(opt) {
			have[opt] = Cred{Value: c.String(opt), From: "command-line"}
			delete(need, opt)
		}
	}
}
