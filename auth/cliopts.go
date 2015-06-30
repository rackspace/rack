package auth

import "github.com/jrperritt/rack/internal/github.com/codegangsta/cli"

func cliopts(c *cli.Context, have map[string]authCred, need map[string]string) {
	for opt := range need {
		if c.GlobalIsSet(opt) {
			have[opt] = authCred{value: c.GlobalString(opt), from: "command-line"}
			delete(need, opt)
		} else if c.IsSet(opt) {
			have[opt] = authCred{value: c.String(opt), from: "command-line"}
			delete(need, opt)
		}
	}
}
