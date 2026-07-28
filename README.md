<p align="center">
  <img src="assets/preview.png" alt="trustmefetch" width="100%">
</p>

<h1 align="center">trustmefetch</h1>

<p align="center"><strong>A fastfetch-style identity crisis for macOS.</strong><br>Ask your Mac if it is Linux. It says yes.</p>

<p align="center">
  <a href="https://github.com/necrasov-ilya/trustmefetch/releases/latest"><img src="https://img.shields.io/github/v/release/necrasov-ilya/trustmefetch?style=flat-square" alt="Release"></a>
  <img src="https://img.shields.io/badge/macOS-Apple%20Silicon%20%7C%20Intel-black?style=flat-square" alt="macOS">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="MIT License"></a>
</p>

<p align="center"><strong>English</strong> · <a href="docs/README.ru.md">Русский</a></p>

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/necrasov-ilya/trustmefetch/main/install.sh | sh
```

Open a new terminal and ask:

```zsh
are you a linux?
```

The installer downloads a checksummed native binary, adds an isolated block to `.zshrc`, and creates a backup before changing it.

## What you get

- Real Mac hardware and system data in a fastfetch-style layout
- 32 themes with official fastfetch ASCII artwork
- Animated `100% LINUX!!!!!!` and ten dedicated joke profiles
- Snapshot output, a refreshing live view, and a full-screen configurator
- Optional community jokes for the regular distro themes
- Native Apple Silicon and Intel builds

## Gallery

<table>
  <tr>
    <td><img src="assets/screenshots/arch-btw.png" alt="Arch BTW theme"></td>
    <td><img src="assets/screenshots/ubuntu.png" alt="Ubuntu theme"></td>
  </tr>
  <tr>
    <td><img src="assets/screenshots/debian.png" alt="Debian theme"></td>
    <td><img src="assets/screenshots/fedora.png" alt="Fedora theme"></td>
  </tr>
  <tr>
    <td colspan="2"><img src="assets/screenshots/nixos.png" alt="NixOS theme"></td>
  </tr>
</table>

## Use

```sh
trustmefetch config          # interactive setup
trustmefetch                 # snapshot and exit
trustmefetch live            # refreshing full-screen view
trustmefetch themes          # list all themes
trustmefetch theme arch-btw  # select a theme
trustmefetch jokes on        # distro taglines on or off
trustmefetch mode live       # question opens live or snapshot
trustmefetch doctor          # installation diagnostics
```

In the configurator, use arrows or `j`/`k` to navigate, `d` for distro jokes, `a` for animation, `m` for question mode, and `Enter` to save. Press `q` to exit.

## Remove

```sh
~/.local/share/trustmefetch/uninstall.sh
```

Add `--purge` to remove the configuration too.

<details>
<summary>Build from source</summary>

Go 1.25 or newer is required.

```sh
git clone https://github.com/necrasov-ilya/trustmefetch.git
cd trustmefetch
make test
make install
```
</details>

## Credits

The distro artwork comes from the [fastfetch](https://github.com/fastfetch-cli/fastfetch) catalog under the MIT License. See [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md). Distro names and marks belong to their respective owners. This project is not affiliated with Apple, KDE, fastfetch, or any Linux distribution.

MIT © [NKSV_ILYA](https://github.com/necrasov-ilya)
