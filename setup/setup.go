package setup

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/jrperritt/rack/internal/github.com/codegangsta/cli"
	"github.com/jrperritt/rack/util"
)

var rackBashAutocomplete = `
#! /bin/bash

_cli_bash_autocomplete() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  if [[ ${prev} != -* ]]; then
    opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi

  if [[ "${#COMP_WORDS[@]}" > 4 ]] && [[ ${cur} != -* ]]; then
    return 0
  fi

  opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
  COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
  return 0
}

complete -F _cli_bash_autocomplete rack
`

// Init runs logic for setting up amenities such as command completion.
func Init(c *cli.Context) {
	w := c.App.Writer
	switch runtime.GOOS {
	case "linux", "darwin":
		rackDir, err := util.RackDir()
		if err != nil {
			fmt.Fprintf(w, "Error running `rack init`: %s\n", err)
			return
		}

		rackCompletionPath := path.Join(rackDir, "bash_autocomplete")
		rackCompletionFile, err := os.Create(rackCompletionPath)
		if err != nil {
			fmt.Fprintf(w, "Error creating `rack` bash completion file: %s\n", err)
			return
		}
		_, err = rackCompletionFile.WriteString(rackBashAutocomplete)
		if err != nil {
			fmt.Fprintf(w, "Error writing to `rack` bash completion file: %s\n", err)
			return
		}
		rackCompletionFile.Close()

		var bashName string
		if runtime.GOOS == "linux" {
			bashName = ".bashrc"
		} else {
			bashName = ".bash_profile"
		}

		homeDir, err := util.HomeDir()
		if err != nil {
			fmt.Fprintf(w, "Unable to access home directory: %s\n", err)
		}

		bashPath := path.Join(homeDir, bashName)
		fmt.Fprintf(w, "Looking for %s in %s\n", bashName, bashPath)
		if _, err := os.Stat(bashPath); os.IsNotExist(err) {
			fmt.Fprintf(w, "%s doesn't exist. You should create it and/or install your operating system's `bash_completion` package.", bashPath)
		} else {
			bashFile, err := os.OpenFile(bashPath, os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintf(w, "Error opening %s: %s\n", bashPath, err)
				return
			}
			defer bashFile.Close()

			sourceContent := fmt.Sprintf("source %s\n", rackCompletionPath)

			bashContentsBytes, err := ioutil.ReadAll(bashFile)
			if strings.Contains(string(bashContentsBytes), sourceContent) {
				fmt.Fprintf(w, "Command completion enabled in %s\n", bashPath)
				return
			}

			_, err = bashFile.WriteString(sourceContent)
			if err != nil {
				fmt.Fprintf(w, "Error writing to %s: %s\n", bashPath, err)
				return
			}

			_, err = exec.Command("/bin/bash", bashPath).Output()
			if err != nil {
				fmt.Fprintf(w, "Error sourcing %s: %s\n", bashPath, err)
				return
			}
			fmt.Fprintf(w, "Command completion enabled in %s\n", bashPath)
			return
		}
	default:
		fmt.Fprintf(w, "Command completion is not currently available for %s\n", runtime.GOOS)
		return
	}
}
