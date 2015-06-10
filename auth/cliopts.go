package auth

import "github.com/codegangsta/cli"

func cliopts(c *cli.Context, have map[string]string, need map[string]string) {
	for opt := range need {
		if c.IsSet(opt) {
			have[opt] = c.String(opt)
			delete(need, opt)
		}
	}
}
