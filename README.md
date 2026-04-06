<p align="center">
  <img src="img/logo.png" alt="VIC logo" width="717">
</p>

  <p align="center">
    <a href="https://github.com/khbarkar/VIC/tags">
      <img src="https://img.shields.io/github/v/tag/khbarkar/VIC?label=tag" alt="Latest Tag">
    </a>
    <a href="https://github.com/khbarkar/VIC/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/khbarkar/VIC" alt="License">
    </a>
    <a href="https://go.dev/">
      <img src="https://img.shields.io/badge/go-1.25.7-00ADD8?logo=go" alt="Go Version">
    </a>
  </p>

# Vanilla Ice Cream 

A small terminal launcher for AI coding CLIs on macOS. 

## Quick Start

```bash
curl -fsSL https://raw.githubusercontent.com/khbarkar/VIC/main/install.sh | bash
```

Then run:

```bash
vic
```

## Update

```bash
vic update
```

## What It Does

- Lets you pick an AI CLI
- Lets you pick a folder
- Opens a new iTerm2 tab
- Splits it into two panes
- Starts the selected CLI on the left
- Opens a shell in the same folder on the right

Supported launchers include Codex, Claude, Gemini, DeepSeek, Grok, Kiro, Cursor, Copilot, and OpenClaw, as long as their CLI command is installed locally.

## Manual Install

```bash
git clone git@github.com:khbarkar/VIC.git
cd vic
make install
vic
```

## Config

VIC reads an optional config file from:

```bash
~/.config/vic/config.json
```

Example:

```json
{
  "initial_folder": "~/private",
  "hidden_projects": ["SnakesInAK8s"]
}
```

If `initial_folder` is not set, VIC defaults to `~/private`.
If `hidden_projects` is set, those direct child folders are hidden from the TUI picker but left untouched on disk.

## License

VIC is licensed under the GNU GPL v3. See [LICENSE](LICENSE).
