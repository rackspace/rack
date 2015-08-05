package setup

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/util"
)

var rackBashAutocomplete = `
#! /bin/bash

_cli_bash_autocomplete() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  # The first 5 words should always be completed by rack
  if [[ ${#COMP_WORDS[@]} -lt 5 ]]; then
    opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
  # All flags should be completed by rack
  elif [[ ${cur} == -* ]]; then
    opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
  # If the previous word wasn't a flag, then the next on has to be, given the 2 conditions above
  elif [[ ${prev} != -* ]]; then
    opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
  fi
  return 0
}

complete -o default -F _cli_bash_autocomplete rack
`

// Init runs logic for setting up amenities such as command completion.
func Init(c *cli.Context) {
	w := c.App.Writer
	rackDir, err := util.RackDir()
	if err != nil {
		fmt.Fprintf(w, "Error running `rack init`: %s\n", err)
		return
	}
	switch runtime.GOOS {
	case "linux", "darwin":
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
	case "windows":
		rackCompletionPath := path.Join(rackDir, "posh_autocomplete.ps1")
		rackCompletionFile, err := os.Create(rackCompletionPath)
		if err != nil {
			fmt.Fprintf(w, "Error creating `rack` PowerShell completion file: %s\n", err)
			return
		}
		_, err = rackCompletionFile.WriteString(rackPoshAutocomplete)
		if err != nil {
			fmt.Fprintf(w, "Error writing to `rack` PowerShell completion file: %s\n", err)
			return
		}
		rackCompletionFile.Close()
	default:
		fmt.Fprintf(w, "Command completion is not currently available for %s\n", runtime.GOOS)
		return
	}
}

var rackPoshAutocomplete = `
function global:TabExpansion2 {
	[CmdletBinding(DefaultParameterSetName = 'ScriptInputSet')]
	Param(
    		[Parameter(ParameterSetName = 'ScriptInputSet', Mandatory = $true, Position = 0)]
    		[string] $inputScript,

    		[Parameter(ParameterSetName = 'ScriptInputSet', Mandatory = $true, Position = 1)]
    		[int] $cursorColumn,

    		[Parameter(ParameterSetName = 'AstInputSet', Mandatory = $true, Position = 0)]
    		[System.Management.Automation.Language.Ast] $ast,

    		[Parameter(ParameterSetName = 'AstInputSet', Mandatory = $true, Position = 1)]
    		[System.Management.Automation.Language.Token[]] $tokens,

    		[Parameter(ParameterSetName = 'AstInputSet', Mandatory = $true, Position = 2)]
    		[System.Management.Automation.Language.IScriptPosition] $positionOfCursor,

    		[Parameter(ParameterSetName = 'ScriptInputSet', Position = 2)]
    		[Parameter(ParameterSetName = 'AstInputSet', Position = 3)]
    		[Hashtable] $options = $null
	)

	End {
    $result = $null

    if ($psCmdlet.ParameterSetName -eq 'ScriptInputSet') {
      $result = [System.Management.Automation.CommandCompletion]::CompleteInput(
        <#inputScript#>  $inputScript,
        <#cursorColumn#> $cursorColumn,
        <#options#>      $options)
    }
    else{
      $result = [System.Management.Automation.CommandCompletion]::CompleteInput(
        <#ast#>              $ast,
        <#tokens#>           $tokens,
        <#positionOfCursor#> $positionOfCursor,
        <#options#>          $options)
    }


    if ($result.CompletionMatches.Count -eq 0){
			if ($psCmdlet.ParameterSetName -eq 'ScriptInputSet') {
        $ast = [System.Management.Automation.Language.Parser]::ParseInput($inputScript, [ref]$tokens, [ref]$null)
      }
      $text = $ast.Extent.Text
    	if($text -match '^*rack.exe*') {
        $cmd1 = $text -split '\s+'
        $end = $cmd1.count - 2
        $cmd2 = $cmd1[0..$end]
        $cmd3 = $cmd2 -join ' '
        $suggestions = Invoke-Expression "$cmd3 --generate-bash-completion"
        ForEach($suggestion in $suggestions) {
          if($suggestion -match $cmd1[$end + 1]) {
            $suggestionObject = New-Object System.Management.Automation.CompletionResult ($suggestion, $suggestion, "Text", $suggestion)
				      $result.CompletionMatches.Add($suggestionObject)
          }
        }
    	}
		}

		return	$result
	}
}
`
