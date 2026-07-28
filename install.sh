#!/bin/sh
set -eu

REPOSITORY="${TRUSTMEFETCH_REPOSITORY:-necrasov-ilya/trustmefetch}"
BIN_DIR="${TRUSTMEFETCH_BIN_DIR:-$HOME/.local/bin}"
SHARE_DIR="${TRUSTMEFETCH_SHARE_DIR:-$HOME/.local/share/trustmefetch}"
ZSHRC="${ZDOTDIR:-$HOME}/.zshrc"

fail() {
  printf 'trustmefetch installer: %s\n' "$1" >&2
  exit 1
}

command -v curl >/dev/null 2>&1 || fail 'curl is required'
command -v tar >/dev/null 2>&1 || fail 'tar is required'

case "$(uname -s)" in
  Darwin) os=Darwin ;;
  *) fail 'only macOS is supported right now' ;;
esac

case "$(uname -m)" in
  arm64) arch=arm64 ;;
  x86_64) arch=x86_64 ;;
  *) fail "unsupported architecture: $(uname -m)" ;;
esac

tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/trustmefetch.XXXXXX")"
trap 'rm -rf "$tmp_dir"' EXIT HUP INT TERM

archive="$tmp_dir/trustmefetch.tar.gz"
url="https://github.com/$REPOSITORY/releases/latest/download/trustmefetch_${os}_${arch}.tar.gz"
printf 'Downloading trustmefetch for %s/%s...\n' "$os" "$arch"
curl -fL --retry 3 --proto '=https' --tlsv1.2 "$url" -o "$archive"
expected_checksum="$(curl -fsSL --retry 3 --proto '=https' --tlsv1.2 "$url.sha256" | awk '{print $1}')"
actual_checksum="$(shasum -a 256 "$archive" | awk '{print $1}')"
[ -n "$expected_checksum" ] || fail 'release checksum is empty'
[ "$expected_checksum" = "$actual_checksum" ] || fail 'release checksum verification failed'
tar -xzf "$archive" -C "$tmp_dir"

test -f "$tmp_dir/trustmefetch" || fail 'release archive does not contain trustmefetch'
test -f "$tmp_dir/trustmefetch.zsh" || fail 'release archive does not contain shell integration'
test -f "$tmp_dir/uninstall.sh" || fail 'release archive does not contain the uninstaller'

mkdir -p "$BIN_DIR" "$SHARE_DIR"
install -m 0755 "$tmp_dir/trustmefetch" "$BIN_DIR/trustmefetch"
install -m 0644 "$tmp_dir/trustmefetch.zsh" "$SHARE_DIR/trustmefetch.zsh"
install -m 0755 "$tmp_dir/uninstall.sh" "$SHARE_DIR/uninstall.sh"

touch "$ZSHRC"
if ! grep -q '^# >>> trustmefetch >>>$' "$ZSHRC"; then
  backup="$ZSHRC.trustmefetch.backup.$(date +%Y%m%d%H%M%S)"
  cp "$ZSHRC" "$backup"
  {
    printf '\n# >>> trustmefetch >>>\n'
    printf 'export PATH="%s:$PATH"\n' "$BIN_DIR"
    printf 'source "%s/trustmefetch.zsh"\n' "$SHARE_DIR"
    printf '# <<< trustmefetch <<<\n'
  } >> "$ZSHRC"
fi

printf '\nInstalled trustmefetch. Open a new terminal, then ask:\n\n'
printf '  you are a linux?\n\n'
printf 'Configure it with: trustmefetch config\n'
