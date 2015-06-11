package output

import (
	"encoding/json"
	"fmt"
)

// JSON prints results in JSON format.
func jsonOut(i interface{}) {
	j, _ := json.Marshal(i)
	fmt.Println(string(j))
}
