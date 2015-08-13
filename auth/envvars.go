package auth

import (
	"os"
	"strings"

	"github.com/rackspace/rack/commandoptions"
)

func envvars(have map[string]commandoptions.Cred, need map[string]string) {
	vars := map[string]string{
		"username": "RS_USERNAME",
		"api-key":  "RS_API_KEY",
		"auth-url": "RS_AUTH_URL",
		"region":   "RS_REGION_NAME",
	}
	for opt := range need {
		if v := os.Getenv(strings.ToUpper(vars[opt])); v != "" {
			have[opt] = commandoptions.Cred{Value: v, From: "environment variable"}
			delete(need, opt)
		}
	}
}
