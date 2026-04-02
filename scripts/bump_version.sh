#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 2 ]; then
  echo "usage: $0 <version-file> <patch|minor|major>" >&2
  exit 1
fi

version_file="$1"
bump_type="$2"

version="$(tr -d '[:space:]' < "$version_file")"
IFS='.' read -r major minor patch <<< "$version"

case "$bump_type" in
  patch)
    patch=$((patch + 1))
    ;;
  minor)
    minor=$((minor + 1))
    patch=0
    ;;
  major)
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  *)
    echo "invalid bump type: $bump_type" >&2
    exit 1
    ;;
esac

printf '%s.%s.%s\n' "$major" "$minor" "$patch" > "$version_file"
