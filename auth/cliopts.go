package auth

import "github.com/codegangsta/cli"

func cliopts(c *cli.Context, have map[string]string, need map[string]string) {
	for opt := range need {
		if c.GlobalIsSet(opt) {
			have[opt] = c.GlobalString(opt)
			delete(need, opt)
		}
	}
}
