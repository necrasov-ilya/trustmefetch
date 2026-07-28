#!/bin/sh
set -eu

root_dir="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
bin_dir="${TRUSTMEFETCH_BIN_DIR:-$HOME/.local/bin}"
share_dir="${TRUSTMEFETCH_SHARE_DIR:-$HOME/.local/share/trustmefetch}"
zshrc="${ZDOTDIR:-$HOME}/.zshrc"

mkdir -p "$bin_dir" "$share_dir"
go build -ldflags "-X main.version=dev" -o "$bin_dir/trustmefetch" "$root_dir/cmd/trustmefetch"
install -m 0644 "$root_dir/shell/trustmefetch.zsh" "$share_dir/trustmefetch.zsh"
install -m 0755 "$root_dir/uninstall.sh" "$share_dir/uninstall.sh"

touch "$zshrc"
if grep -q '^# A tiny identity crisis for macOS\.$' "$zshrc"; then
  printf 'An older manual trustmefetch snippet exists in %s. Remove it before enabling the managed integration.\n' "$zshrc" >&2
elif ! grep -q '^# >>> trustmefetch >>>$' "$zshrc"; then
  cp "$zshrc" "$zshrc.trustmefetch.backup"
  {
    printf '\n# >>> trustmefetch >>>\n'
    printf 'export PATH="%s:$PATH"\n' "$bin_dir"
    printf 'source "%s/trustmefetch.zsh"\n' "$share_dir"
    printf '# <<< trustmefetch <<<\n'
  } >> "$zshrc"
fi

printf 'Local build installed to %s/trustmefetch\n' "$bin_dir"
