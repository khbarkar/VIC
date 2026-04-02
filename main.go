package main

import (
	"ai/config"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ai/launcher"
	"ai/ui"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	cli         launcher.CLIInfo
	kind        string
	label       string
	description string
	path        string
	installed   bool
}

func (i item) Title() string {
	switch i.kind {
	case "action":
		return ui.ActionStyle.Render(i.label)
	case "dir":
		return ui.ProjectStyle.Render(i.label)
	}

	color := ui.GetCLIColor(i.cli.ID)
	name := lipgloss.NewStyle().Foreground(color).Bold(true).Render(i.cli.Name)
	status := ui.InstalledStyle.Render("[installed]")
	if !i.installed {
		status = ui.MissingStyle.Render("[missing]")
	}
	return fmt.Sprintf("%-20s %s", name, status)
}

func (i item) Description() string {
	if i.description != "" {
		return i.description
	}
	return "Launch " + i.cli.Name
}

func (i item) FilterValue() string {
	if i.kind == "cli" {
		return i.cli.Name
	}
	return i.label
}

type model struct {
	cliList     list.Model
	projectList list.Model
	nameInput   textinput.Model
	state       string
	selectedCLI launcher.CLIInfo
	privateDir  string
	err         error
	winWidth    int
	winHeight   int
}

func initialModel() model {
	home, _ := os.UserHomeDir()
	privateDir := filepath.Join(home, "private")
	cfg, _ := config.LoadConfig()
	if cfg != nil && strings.TrimSpace(cfg.InitialFolder) != "" {
		if expanded, err := launcher.ExpandPath(cfg.InitialFolder); err == nil {
			privateDir = expanded
		}
	}
	_ = os.MkdirAll(privateDir, 0755)

	cliList := newList(buildCLIItems(), "Choose AI")
	projectList := newList(buildProjectItems(privateDir), "Choose Project")

	nameInput := textinput.New()
	nameInput.Placeholder = "project-name"
	nameInput.CharLimit = 64
	nameInput.Width = 36

	return model{
		cliList:     cliList,
		projectList: projectList,
		nameInput:   nameInput,
		state:       "list",
		privateDir:  privateDir,
	}
}

func newList(items []list.Item, title string) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.Styles.Title = ui.ListTitleStyle
	l.Styles.NoItems = ui.EmptyStateStyle
	return l
}

func buildCLIItems() []list.Item {
	items := []list.Item{}
	for _, cli := range launcher.CLIs {
		_, err := launcher.ResolveCommand(cli)
		items = append(items, item{
			cli:       cli,
			kind:      "cli",
			installed: err == nil,
		})
	}
	return items
}

func buildProjectItems(privateDir string) []list.Item {
	entries, err := os.ReadDir(privateDir)
	if err != nil {
		return []list.Item{
			item{kind: "action", label: "+ Create new project", description: "Create a folder directly under ~/private"},
		}
	}

	items := make([]list.Item, 0, len(entries)+1)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		items = append(items, item{
			kind:        "dir",
			label:       name,
			description: filepath.Join("~", "private", name),
			path:        filepath.Join(privateDir, name),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		left := items[i].(item)
		right := items[j].(item)
		return strings.ToLower(left.label) < strings.ToLower(right.label)
	})

	items = append(items, item{
		kind:        "action",
		label:       "+ Create new project",
		description: "Create a folder directly under ~/private",
	})

	return items
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height
		listWidth := min(72, msg.Width-8)
		listHeight := max(8, msg.Height-12)
		m.cliList.SetSize(listWidth, listHeight)
		m.projectList.SetSize(listWidth, listHeight)
		m.nameInput.Width = min(40, max(20, msg.Width-20))

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.err != nil {
				m.err = nil
				return m, nil
			}
			if m.state == "list" {
				return m, tea.Quit
			}
			if m.state == "create_dir" {
				m.state = "dir_pick"
				m.nameInput.Blur()
				return m, nil
			}
			m.state = "list"
			return m, nil
		case "enter":
			if m.state == "list" {
				i, ok := m.cliList.SelectedItem().(item)
				if !ok {
					break
				}
				m.selectedCLI = i.cli
				m.state = "dir_pick"
				m.projectList.SetItems(buildProjectItems(m.privateDir))
				return m, nil
			} else if m.state == "dir_pick" {
				i, ok := m.projectList.SelectedItem().(item)
				if !ok {
					break
				}
				if i.kind == "action" {
					m.state = "create_dir"
					m.nameInput.SetValue("")
					m.nameInput.Focus()
					return m, nil
				}
				return m.launchWith(i.path)
			} else if m.state == "create_dir" {
				return m.createAndLaunch()
			}
		}

	case errorMsg:
		m.err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	switch m.state {
	case "list":
		m.cliList, cmd = m.cliList.Update(msg)
	case "dir_pick":
		m.projectList, cmd = m.projectList.Update(msg)
	case "create_dir":
		m.nameInput, cmd = m.nameInput.Update(msg)
	}

	return m, cmd
}

func (m model) createAndLaunch() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.nameInput.Value())
	if name == "" {
		m.err = fmt.Errorf("project name cannot be empty")
		return m, nil
	}
	if strings.Contains(name, "/") {
		m.err = fmt.Errorf("project name must be a single directory name")
		return m, nil
	}

	dir := filepath.Join(m.privateDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		m.err = err
		return m, nil
	}

	return m.launchWith(dir)
}

func (m model) launchWith(dir string) (tea.Model, tea.Cmd) {
	err := launcher.Launch(m.selectedCLI, dir)
	if err != nil {
		m.err = err
		return m, nil
	}
	return m, tea.Quit
}

type errorMsg struct{ err error }

func (m model) View() string {
	if m.err != nil {
		errText := ui.MissingStyle.Render("Error: " + m.err.Error())
		help := ui.SubTitleStyle.Render("Press q to go back")
		return ui.AppFrameStyle.Copy().
			Width(max(1, m.winWidth)).
			Height(max(1, m.winHeight)).
			Render(lipgloss.JoinVertical(lipgloss.Left, ui.AlertStyle.Render(errText), help))
	}

	cardHeight := 0
	cardWidth := equalPanelWidth(m.winWidth)
	subtitle := ui.SubTitleStyle.Render("Choose one model, choose one folder under ~/private, and launch into a split iTerm2 tab.")

	switch m.state {
	case "list":
		body := renderTwoPanelLayout(
			m.winWidth,
			ui.PillStyle.Render("CHOOSE MODEL"),
			renderSelectablePanel(m.cliList, cardWidth, cardHeight, renderCLIItemBody),
			ui.PillStyle.Render("DETAIL"),
			renderDetailPanel(selectedListItem(m.cliList), "Select a model to continue.", renderCLISelectionDetail, cardWidth, cardHeight),
		)
		help := ui.RenderHelp("up/down move", "enter continue", "q quit")
		return renderAppFrame(m.winWidth, m.winHeight, subtitle, body, help)

	case "dir_pick":
		body := renderTwoPanelLayout(
			m.winWidth,
			ui.PillStyle.Render("CHOOSE PROJECT"),
			renderSelectablePanel(m.projectList, cardWidth, cardHeight, renderProjectItemBody),
			ui.PillStyle.Render("LAUNCH"),
			renderInfoPanel([]string{
				"AI: " + m.selectedCLI.Name,
				"Only direct children of ~/private are shown.",
				"Choose one folder or create a new project.",
			}, cardWidth, cardHeight),
		)
		help := ui.RenderHelp("up/down move", "enter launch", "q back")
		return renderAppFrame(m.winWidth, m.winHeight, subtitle, body, help)

	case "create_dir":
		form := ui.RenderPanel(lipgloss.JoinVertical(
			lipgloss.Left,
			ui.PillStyle.Render("NEW PROJECT"),
			"",
			ui.SubTitleStyle.Render("Create one folder directly under ~/private."),
			"",
			ui.InputShellStyle.Render(m.nameInput.View()),
		), true, cardWidth, cardHeight)
		body := renderTwoPanelLayout(
			m.winWidth,
			"",
			form,
			"",
			renderInfoPanel([]string{
				"Model: " + m.selectedCLI.Name,
				"Folder must be a single directory name.",
				"Launch happens immediately after create.",
			}, cardWidth, cardHeight),
		)
		help := ui.RenderHelp("type folder name", "enter create and launch", "q cancel")
		return renderAppFrame(m.winWidth, m.winHeight, subtitle, body, help)
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func renderTwoPanelLayout(width int, leftLabel, leftBody, rightLabel, rightBody string) string {
	left := leftBody
	if leftLabel != "" {
		left = lipgloss.JoinVertical(lipgloss.Left, leftLabel, "", leftBody)
	}
	right := rightBody
	if rightLabel != "" {
		right = lipgloss.JoinVertical(lipgloss.Left, rightLabel, "", rightBody)
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		lipgloss.NewStyle().MarginLeft(3).Render(right),
	)

	return lipgloss.NewStyle().
		Background(ui.BgBottomColor).
		Width(max(60, width-4)).
		AlignHorizontal(lipgloss.Center).
		Render(content)
}

func renderAppFrame(width, height int, parts ...string) string {
	body := lipgloss.JoinVertical(lipgloss.Center, parts...)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Background(ui.BgBottomColor).
			Width(max(1, width-4)).
			AlignHorizontal(lipgloss.Center).
			Render(body),
	)
	return ui.AppFrameStyle.Copy().
		Width(max(1, width)).
		Height(max(1, height)).
		Render(content)
}

func equalPanelWidth(width int) int {
	usable := min(110, max(72, width-10))
	return max(30, (usable-3)/2)
}

func renderSelectablePanel(l list.Model, width, height int, renderer func(item, bool) string) string {
	items := l.Items()
	if len(items) == 0 {
		return ui.RenderPanel(ui.EmptyStateStyle.Render("Nothing here yet."), true, width, height)
	}

	lines := make([]string, 0, len(items))
	selected := l.Index()
	for idx, raw := range items {
		entry, ok := raw.(item)
		if !ok {
			continue
		}
		lines = append(lines, renderer(entry, idx == selected))
	}

	return ui.RenderPanel(lipgloss.JoinVertical(lipgloss.Left, lines...), true, width, height)
}

func renderDetailPanel(selected *item, empty string, renderer func(item) string, width, height int) string {
	content := ui.EmptyStateStyle.Render(empty)
	if selected != nil {
		content = renderer(*selected)
	}
	return ui.RenderPanel(content, false, width, height)
}

func renderInfoPanel(lines []string, width, height int) string {
	rendered := make([]string, 0, len(lines))
	for _, line := range lines {
		rendered = append(rendered, ui.ItemBodyStyle.Render(line))
	}
	return ui.RenderPanel(lipgloss.JoinVertical(lipgloss.Left, rendered...), false, width, height)
}

func renderCLIItemBody(entry item, selected bool) string {
	nameLine := "  " + entry.cli.Name
	nameStyle := ui.ProjectStyle
	if selected {
		nameLine = "> " + entry.cli.Name
		nameStyle = ui.SelectedItemStyle
	}

	status := "installed"
	if !entry.installed {
		status = "missing"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		nameStyle.Render(nameLine),
		ui.ItemMetaStyle.Render("  "+status),
	)
}

func renderProjectItemBody(entry item, selected bool) string {
	titleLine := "  " + entry.label
	titleStyle := ui.ProjectStyle
	if entry.kind == "action" {
		titleStyle = ui.ActionStyle
	}
	if selected {
		titleLine = "> " + entry.label
		titleStyle = ui.SelectedItemStyle
	}

	meta := entry.description
	if entry.kind == "action" {
		meta = "Create a new direct child under ~/private"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(titleLine),
		ui.ItemMetaStyle.Render("  "+meta),
	)
}

func renderCLISelectionDetail(entry item) string {
	status := "Installed"
	if !entry.installed {
		status = "Missing"
	}

	lines := []string{
		ui.ProjectStyle.Render(entry.cli.Name),
		ui.MetaStyle.Render(status + "  " + entry.cli.Category),
	}
	if entry.cli.About != "" {
		lines = append(lines, "", ui.ItemBodyStyle.Render(entry.cli.About))
	}
	if len(entry.cli.Models) > 0 {
		lines = append(lines, "", ui.ItemMetaStyle.Render("Available models"))
		for _, model := range entry.cli.Models {
			lines = append(lines, ui.ItemBodyStyle.Render("• "+model))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lines...,
	)
}

func renderProjectSelectionDetail(entry item) string {
	if entry.kind == "action" {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			ui.ActionStyle.Render(entry.label),
			"",
			ui.ItemBodyStyle.Render("Creates a new folder directly under ~/private and launches the selected AI there."),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		ui.ProjectStyle.Render(entry.label),
		ui.MetaStyle.Render(entry.description),
		"",
		ui.ItemBodyStyle.Render("Ready to launch in this project."),
	)
}

func selectedListItem(l list.Model) *item {
	rawItems := l.Items()
	if len(rawItems) == 0 {
		return nil
	}

	index := l.Index()
	if index < 0 || index >= len(rawItems) {
		index = 0
	}

	entry, ok := rawItems[index].(item)
	if !ok {
		return nil
	}
	return &entry
}
