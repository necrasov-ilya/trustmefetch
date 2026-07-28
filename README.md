# trustmefetch

**A fastfetch-style identity crisis for macOS. Ask if it is Linux. It says yes.**

`trustmefetch` collects real information from your Mac, dresses it as one of 32 Linux personalities, and answers the only question that matters:

```console
$ are you a linux?

                  -`                     ilya@MacBook-Pro
                 .o+`                    -----------------
                `ooo/                    OS: Arch Linux arm64
               `+oooo:                   Host: MacBookPro18,2
              `+oooooo:                  Kernel: Darwin 25.5.0
              -+oooooo+:                 Desktop: KDE Plasma 6
            `/:-:++oooo+:                Verdict: yes. trust me.
```

The host, OS build, kernel, uptime, package count, shell, terminal, display, CPU, GPU, memory, disk, and battery values come from the current machine. The selected distro, desktop, logo, palette, and questionable confidence are the joke.

## Features

- 32 built-in themes
- Ten joke profiles, including animated `100% LINUX!!!!!!`
- 22 distro profiles using the official fastfetch ASCII logo catalog
- Full-screen theme picker with live preview
- Full-screen live view with periodically refreshed CPU, memory, disk, battery, and uptime values
- Real macOS system discovery without hardcoded hardware values
- True-color output with a readable no-color fallback
- Exact zsh phrase `are you a linux?`
- Native Apple Silicon and Intel binaries
- Idempotent installer and removable shell integration

## Install

Once the first release is published:

```sh
curl -fsSL https://raw.githubusercontent.com/necrasov-ilya/trustmefetch/main/install.sh | sh
```

Open a new terminal and ask:

```zsh
are you a linux?
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
trustmefetch jokes on
trustmefetch jokes off
trustmefetch doctor
```

Inside the picker, use arrow keys or `j` and `k` to move, `Enter` or `s` to save, `r` for a random theme, `a` to toggle animation, `m` to choose the question mode, `d` to toggle taglines for distro themes, and `q` to quit. Joke themes always keep their taglines.

## Snapshot and live modes

Like fastfetch, the regular command prints a snapshot and returns to the shell:

```sh
trustmefetch
```

The live view stays open and refreshes system values until `q` is pressed:

```sh
trustmefetch live
```

Choose what the question opens:

```sh
trustmefetch mode snapshot
trustmefetch mode live
```

The same setting is available inside `trustmefetch config` with the `m` key. In live mode, use `Space` to pause, `r` to refresh, `t` to switch themes, arrow keys to scroll, and `q` to exit.

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

The distro names and marks belong to their respective owners. This project is not affiliated with Apple, KDE, fastfetch, or any Linux distribution. The bundled distro ASCII artwork comes from fastfetch under the MIT License. See [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

## License

MIT. See [LICENSE](LICENSE).
