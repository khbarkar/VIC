<p align="center">
  <img src="img/logo.png" alt="VIC logo" width="320">
</p>

# VIC

Vanilla Ice Cream.

A small terminal launcher for AI coding CLIs on macOS.

## Quick Start

```bash
git clone git@github.com:khbarkar/vic.git
cd vic
make install
vic
```

## Config

VIC reads an optional config file from:

```bash
~/.config/ai-launcher/config.json
```

Example:

```json
{
  "initial_folder": "~/private"
}
```

If `initial_folder` is not set, VIC defaults to `~/private`.

## What It Does

- Lets you pick an AI CLI
- Lets you pick a folder
- Opens a new iTerm2 tab
- Splits it into two panes
- Starts the selected CLI on the left
- Opens a shell in the same folder on the right
