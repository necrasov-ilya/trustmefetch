# Shell integration for trustmefetch.
# This file is sourced from ~/.zshrc by the installer.

function _trustmefetch_you() {
  if [[ "$*" == 'are a linux?' ]]; then
    command trustmefetch --question
    return $?
  fi

  print -u2 -P '%F{yellow}Try:%f you are a linux?'
  return 2
}

# `noglob` keeps the question mark literal without changing global zsh options.
alias you='noglob _trustmefetch_you'

