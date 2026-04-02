# VIC Launcher

<p align="center">
  <img src="./logo.png" alt="VIC logo" width="320" />
</p>

A TUI launcher for multiple AI CLI tools built with Go and Bubbletea.

## Features

- Launch AI CLIs (Gemini, Codex, Kiro, Cursor, Grok, Copilot) in split iTerm2 panes
- Directory picker with autocomplete navigation
- Install/update CLIs from the maintenance menu
- Reads configuration from `~/.config/ai-launcher/config.json`

## Installation

```bash
cd ~/private/ai
make install
```

This builds the binary and symlinks `~/.local/bin/vic` → `~/private/ai/vic`.

## Usage

```bash
vic
```

Optional startup config:

```json
{
  "initial_folder": "~/private"
}
```

If `initial_folder` is not set, `vic` defaults to `~/private`.

### Navigation

**Main Menu:**
- `↑/↓` - Select CLI
- `Enter` - Launch selected CLI
- `q` - Quit

**Directory Picker:**
- `↑/↓` - Navigate entries
- `Enter` - Navigate into subdirectory
- `Space` - Confirm current directory and launch
- `q` - Back to main menu

**Maintenance:**
- Select a CLI to install or update it

## Development

```bash
make build    # Build binary
make run      # Build and run
make test     # Run tests
make lint     # Run linter
make clean    # Remove binary
```

## Supported CLIs

- **Gemini** - `gemini`
- **Codex** - `codex`
- **Kiro** - `kiro-cli`
- **Cursor** - `agent` / `cursor-agent`
- **Grok** - `grok-cli` / `grok`
- **Copilot** - `copilot-cli` / `copilot`

## How It Works

When you launch a CLI:
1. Opens a new tab in your current iTerm2 window
2. Splits the tab vertically
3. Left pane: runs the AI CLI in your selected directory
4. Right pane: shell in the same directory
5. Sets cursor color and badge for visual identification
