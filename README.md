# trustmefetch

**A fastfetch-style identity crisis for macOS. Ask if it is Linux. It says yes.**

`trustmefetch` collects real information from your Mac, dresses it as one of 32 Linux personalities, and answers the only question that matters:

```console
$ you are a linux?

       .--.        ilya@MacBook-Pro
      |o_o |       ─────────────────
      |:_/ |       OS         : TrustMe Linux arm64
     //   \ \      Host       : MacBookPro18,2
    (|     | )     Kernel     : Darwin 25.5.0
   /'\_   _/\      Desktop    : KDE Plasma 6
   \___)=(___/     Verdict    : yes. trust me.

                   Trust me, it's Linux.
```

The host, OS build, kernel, uptime, package count, shell, terminal, display, CPU, GPU, memory, disk, and battery values come from the current machine. The selected distro, desktop, logo, palette, and questionable confidence are the joke.

## Features

- 32 built-in themes with original terminal-sized ASCII artwork
- Ten joke profiles, including animated `100% LINUX!!!!!!`
- 22 familiar distro profiles inspired by conventional fetch tools
- Full-screen theme picker with live preview
- Real macOS system discovery without hardcoded hardware values
- True-color output with a readable no-color fallback
- Exact zsh phrase `you are a linux?`
- Native Apple Silicon and Intel binaries
- Idempotent installer and removable shell integration

## Install

Once the first release is published:

```sh
curl -fsSL https://raw.githubusercontent.com/necrasov-ilya/trustmefetch/main/install.sh | sh
```

Open a new terminal and ask:

```zsh
you are a linux?
```

To inspect the installer before running it:

```sh
curl -fsSL https://raw.githubusercontent.com/necrasov-ilya/trustmefetch/main/install.sh -o install.sh
less install.sh
sh install.sh
```

## Configure

Open the interactive picker:

```sh
trustmefetch config
```

Or configure it directly:

```sh
trustmefetch themes
trustmefetch theme arch-btw
trustmefetch preview rgb-linux
trustmefetch random
trustmefetch doctor
```

Inside the picker, use arrow keys or `j` and `k` to move, `Enter` or `s` to save, `r` for a random theme, `a` to toggle animation, and `q` to quit.

## Build locally

Go 1.25 or newer is required.

```sh
git clone https://github.com/necrasov-ilya/trustmefetch.git
cd trustmefetch
make test
make install
```

The local installer places the binary in `~/.local/bin`, shell integration in `~/.local/share/trustmefetch`, and configuration in `~/.config/trustmefetch/config.json`.

## Uninstall

Use the uninstaller placed by the installer:

```sh
~/.local/share/trustmefetch/uninstall.sh
```

Keep the configuration by default or remove it too:

```sh
~/.local/share/trustmefetch/uninstall.sh --purge
```

The installer creates a backup before modifying `.zshrc`. The uninstaller also makes a backup before removing the managed block.

## Themes

The catalog includes `rgb-linux`, `seriously-linux`, `trust-me-bro`, `arch-btw`, `kde-delusion`, `sudo-believe`, `kernel-of-truth`, `macos-who`, `penguin-certified`, `retina-wayland`, Arch, Ubuntu, Debian, Fedora, NixOS, Alpine, Manjaro, EndeavourOS, openSUSE, Gentoo, Kali, Linux Mint, Pop!_OS, elementary OS, Void, Slackware, Rocky Linux, AlmaLinux, CentOS Stream, Garuda, CachyOS, and KDE neon.

The distro names and marks belong to their respective owners. This project is not affiliated with Apple, KDE, fastfetch, or any Linux distribution.

## License

MIT. See [LICENSE](LICENSE).
