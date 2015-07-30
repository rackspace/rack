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
