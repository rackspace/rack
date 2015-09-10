package stacks

import (
	"errors"
	"fmt"
	"strings"
)

type Template struct {
	TE
}

var TemplateFormatVersions = map[string]bool{
	"HeatTemplateFormatVersion": true,
	"heat_template_version":     true,
	"AWSTemplateFormatVersion":  true,
}

func (t *Template) Validate() error {
	if t.Parsed == nil {
		if err := t.Parse(); err != nil {
			return err
		}
	}
	for key, _ := range t.Parsed {
		if _, ok := TemplateFormatVersions[key]; ok {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Template format version not found."))
}

// function to choose keys whose values are other template files
func ignoreIfTemplate(key string, value interface{}) bool {
	if key != "get_file" && key != "type" {
		return true
	}
	valueString, ok := value.(string)
	if !ok {
		return true
	}
	if key == "type" && !(strings.HasSuffix(valueString, ".template") || strings.HasSuffix(valueString, ".yaml")) {
		return true
	}
	return false
}
