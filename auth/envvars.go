package auth

import (
	"os"
	"strings"
)

func envvars(have map[string]authCred, need map[string]string) {
	vars := map[string]string{
		"username": "RS_USERNAME",
		"apikey":   "RS_API_KEY",
		"authurl":  "RS_AUTH_URL",
		"region":   "RS_REGION_NAME",
	}
	for opt := range need {
		if v := os.Getenv(strings.ToUpper(vars[opt])); v != "" {
			have[opt] = authCred{value: v, from: "environment variable"}
			delete(need, opt)
		}
	}
}
