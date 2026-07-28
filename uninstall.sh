#!/bin/sh
set -eu

BIN_DIR="${TRUSTMEFETCH_BIN_DIR:-$HOME/.local/bin}"
SHARE_DIR="${TRUSTMEFETCH_SHARE_DIR:-$HOME/.local/share/trustmefetch}"
ZSHRC="${ZDOTDIR:-$HOME}/.zshrc"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/trustmefetch"

if [ -f "$ZSHRC" ]; then
  tmp_file="$(mktemp "${TMPDIR:-/tmp}/trustmefetch-zshrc.XXXXXX")"
  awk '
    $0 == "# >>> trustmefetch >>>" { skip = 1; next }
    $0 == "# <<< trustmefetch <<<" { skip = 0; next }
    !skip { print }
  ' "$ZSHRC" > "$tmp_file"
  cp "$ZSHRC" "$ZSHRC.trustmefetch.uninstall-backup"
  mv "$tmp_file" "$ZSHRC"
fi

rm -f "$BIN_DIR/trustmefetch"
rm -f "$SHARE_DIR/trustmefetch.zsh"
rm -f "$SHARE_DIR/uninstall.sh"
rmdir "$SHARE_DIR" 2>/dev/null || true

if [ "${1:-}" = "--purge" ]; then
  rm -f "$CONFIG_DIR/config.json"
  rmdir "$CONFIG_DIR" 2>/dev/null || true
  printf 'Removed trustmefetch and its configuration.\n'
else
  printf 'Removed trustmefetch. Configuration was kept in ~/.config/trustmefetch.\n'
fi
