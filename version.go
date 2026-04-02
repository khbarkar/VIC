package main

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var embeddedVersion string

var (
	Version = strings.TrimSpace(embeddedVersion)
	Commit  = "unknown"
)
