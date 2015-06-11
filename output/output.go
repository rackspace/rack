package output

import "github.com/codegangsta/cli"

// Print prints the results of the CLI command.
func Print(c *cli.Context, f *func() interface{}, keys []string) {
	i := (*f)()
	if c.IsSet("json") {
		jsonOut(i)
		return
	}
	if len(keys) == 0 {
		(*f)()
		return
	}
	if c.IsSet("csv") {
		switch i.(type) {
		case map[string]interface{}:
			m := i.(map[string]interface{})
			metadataCSV(c, m, keys)
		case []map[string]interface{}:
			m := i.([]map[string]interface{})
			listCSV(c, m, keys)
		}

		return
	}
	switch i.(type) {
	case map[string]interface{}:
		m := i.(map[string]interface{})
		metadataTable(c, m, keys)
	case []map[string]interface{}:
		m := i.([]map[string]interface{})
		listTable(c, m, keys)
	}
}
