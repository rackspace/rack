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
    COMPREPLY=()
    return 0
  fi

  opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
  COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
  return 0
}

complete -o default -F _cli_bash_autocomplete rack
