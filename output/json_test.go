package output

import (
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
)

func TestLimitJSONFields(t *testing.T) {
	m := map[string]interface{}{"id": "12345", "name": "name", "status": "actvive"}
	keys := []string{"id", "status"}
	expected := map[string]interface{}{"id": "12345", "status": "actvive"}
	actual := limitJSONFields(m, keys)
	th.AssertDeepEquals(t, expected, actual)

	keys = []string{}
	expected = m
	actual = limitJSONFields(m, keys)
	th.AssertDeepEquals(t, expected, actual)

	keys = []string{"not", "there"}
	expected = map[string]interface{}{}
	actual = limitJSONFields(m, keys)
	th.AssertDeepEquals(t, expected, actual)
}
