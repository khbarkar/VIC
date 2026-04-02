package launcher

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CLIInfo struct {
	ID          string
	Name        string
	PrimaryCmd  string
	FallbackCmd string
	CursorColor string
	Category    string
	About       string
	Models      []string
}

var CLIs = []CLIInfo{
	{ID: "gemini", Name: "Gemini", PrimaryCmd: "gemini", CursorColor: "#4A90E2", Category: "Code AI", About: "Google CLI for Gemini chat and coding workflows.", Models: []string{"Gemini 2.5 Pro", "Gemini 2.5 Flash"}},
	{ID: "codex", Name: "Codex", PrimaryCmd: "codex", CursorColor: "#8B5E3C", Category: "Code AI", About: "OpenAI coding agent in the terminal.", Models: []string{"GPT-5", "GPT-5 Mini"}},
	{ID: "kiro", Name: "Kiro", PrimaryCmd: "kiro-cli", CursorColor: "#8A2BE2", Category: "Code AI", About: "Kiro terminal chat workflow.", Models: []string{"Kiro Chat"}},
	{ID: "grok", Name: "Grok", PrimaryCmd: "grok-cli", FallbackCmd: "grok", CursorColor: "#D4A017", Category: "Code AI", About: "xAI terminal assistant.", Models: []string{"Grok 4"}},
	{ID: "claude", Name: "Claude", PrimaryCmd: "claude", CursorColor: "#8B4513", Category: "Code AI", About: "Anthropic CLI for chat and code tasks.", Models: []string{"Claude Sonnet", "Claude Opus"}},
	{ID: "cursor", Name: "Cursor", PrimaryCmd: "agent", FallbackCmd: "cursor-agent", CursorColor: "#2E8B57", Category: "Tool", About: "Cursor agent workflow from the terminal.", Models: []string{"Configured in Cursor"}},
	{ID: "copilot", Name: "Copilot", PrimaryCmd: "copilot-cli", FallbackCmd: "copilot", CursorColor: "#FF5FA2", Category: "Tool", About: "GitHub Copilot command-line workflow.", Models: []string{"Configured in Copilot"}},
	{ID: "openclaw", Name: "OpenClaw", PrimaryCmd: "openclaw", CursorColor: "#DC143C", Category: "Tool", About: "OpenClaw terminal workflow.", Models: []string{"Configured in OpenClaw"}},
}

func escapeAppleScript(s string) string {
	// Escape backslashes first, then double quotes
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func ResolveCommand(cli CLIInfo) (string, error) {
	path, err := exec.LookPath(cli.PrimaryCmd)
	if err == nil {
		return path, nil
	}
	if cli.FallbackCmd != "" {
		path, err = exec.LookPath(cli.FallbackCmd)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("command not found for %s", cli.Name)
}

func resolveLaunchArgs(cli CLIInfo, cmdPath string) ([]string, error) {
	args := []string{cmdPath}

	if usesEnvNode(cmdPath) {
		nodePath, err := exec.LookPath("node")
		if err != nil {
			return nil, fmt.Errorf("node is required to launch %s: %w", cli.Name, err)
		}
		args = []string{nodePath, cmdPath}
	}

	switch cli.ID {
	case "kiro":
		args = append(args, "chat")
	}

	return args, nil
}

func usesEnvNode(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return false
	}

	return strings.TrimSpace(scanner.Text()) == "#!/usr/bin/env node"
}

func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return filepath.Abs(path)
}

func SessionBootstrap(_badge, color, title string) string {

	return fmt.Sprintf(
		`echo -e '\033]12;%s\a\033]0;%s\a'`,
		color, title,
	)
}

func Launch(cli CLIInfo, dir string) error {
	cmdPath, err := ResolveCommand(cli)
	if err != nil {
		return err
	}

	absDir, err := ExpandPath(dir)
	if err != nil {
		return err
	}

	badge := cli.Name
	title := badge

	cliArgs, err := resolveLaunchArgs(cli, cmdPath)
	if err != nil {
		return err
	}

	quotedCliArgs := make([]string, len(cliArgs))
	for i, arg := range cliArgs {
		quotedCliArgs[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(arg, "'", "'\\''"))
	}
	cliCmd := strings.Join(quotedCliArgs, " ")

	cliLine := fmt.Sprintf("cd '%s'; %s; %s", strings.ReplaceAll(absDir, "'", "'\\''"), SessionBootstrap(badge, cli.CursorColor, title), cliCmd)
	dirLine := fmt.Sprintf("cd '%s'", strings.ReplaceAll(absDir, "'", "'\\''"))

	applescript := fmt.Sprintf(`
on run argv
  set cliCommand to "%s"
  set dirCommand to "%s"

  tell application "iTerm2"
    activate
    if (count windows) = 0 then
      set myWindow to (create window with default profile)
      set newTab to current tab of myWindow
    else
      set myWindow to current window
      tell myWindow
        set newTab to (create tab with default profile)
      end tell
    end if

    set primarySession to current session of newTab
    set secondarySession to (split vertically with default profile) of primarySession

    tell primarySession
      write text dirCommand
    end tell

    tell secondarySession
      write text cliCommand
    end tell
  end tell
end run
`, escapeAppleScript(cliLine), escapeAppleScript(dirLine))

	cmd := exec.Command("osascript", "-")
	cmd.Stdin = strings.NewReader(applescript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("osascript failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func Update(repoDir string) error {
	absRepoDir, err := ExpandPath(repoDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(absRepoDir, ".git")); err != nil {
		return fmt.Errorf("no git checkout found at %s", absRepoDir)
	}

	pullCmd := exec.Command("git", "-C", absRepoDir, "pull", "--ff-only")
	pullOutput, err := pullCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %w\n%s", err, strings.TrimSpace(string(pullOutput)))
	}

	installCmd := exec.Command("make", "-C", absRepoDir, "install")
	installOutput, err := installCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("make install failed: %w\n%s", err, strings.TrimSpace(string(installOutput)))
	}

	return nil
}
