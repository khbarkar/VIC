<p align="center">
  <img src="img/logo.png" alt="VIC logo" width="717">
</p>

<p align="center">
  <a href="https://github.com/khbarkar/vic/releases/latest">
    <img src="https://img.shields.io/github/v/release/khbarkar/vic?label=release" alt="Latest Release">
  </a>
  <a href="https://github.com/khbarkar/vic/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/khbarkar/vic" alt="License">
  </a>
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/go-1.25.7-00ADD8?logo=go" alt="Go Version">
  </a>
  <a href="https://github.com/khbarkar/vic/actions/workflows/release.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/khbarkar/vic/release.yml?label=release%20workflow" alt="Release Workflow">
  </a>
</p>

# Vanilla Ice Cream 

A small terminal launcher for AI coding CLIs on macOS. 

## Quick Start

```bash
curl -fsSL https://raw.githubusercontent.com/khbarkar/vic/main/install.sh | bash
```

Then run:

```bash
vic
```

## Update

```bash
vic update
```

If you are running from a local checkout:

```bash
make update
```

## Release Versioning

VIC uses semantic versioning:

- patch: bug fixes
- minor: new features
- major: breaking changes

The current version is stored in `VERSION`.

On pushes to `main`, the release workflow bumps and tags automatically:

- `fix:` or any non-feature change -> patch
- `feat:` -> minor
- `BREAKING CHANGE` or `!:` -> major
## What It Does

- Lets you pick an AI CLI
- Lets you pick a folder
- Opens a new iTerm2 tab
- Splits it into two panes
- Starts the selected CLI on the left
- Opens a shell in the same folder on the right

## Manual Install

```bash
git clone git@github.com:khbarkar/vic.git
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
