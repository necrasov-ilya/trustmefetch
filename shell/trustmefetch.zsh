# Shell integration for trustmefetch.
# This file is sourced from ~/.zshrc by the installer.

function _trustmefetch_are() {
	if [[ "$*" == 'you a linux?' ]]; then
		command trustmefetch --question
		return $?
	fi

	print -u2 -P '%F{yellow}Try:%f are you a linux?'
	return 2
}

# `noglob` keeps the question mark literal without changing global zsh options.
alias are='noglob _trustmefetch_are'
