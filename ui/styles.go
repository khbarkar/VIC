package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	BgTopColor    = lipgloss.Color("#171A1F")
	BgBottomColor = lipgloss.Color("#171A1F")
	CardBgColor   = lipgloss.Color("#171A1F")
	TextColor     = lipgloss.Color("#F2E7D8")
	MutedColor    = lipgloss.Color("#C7B7A5")
	BlueColor     = lipgloss.Color("#4E8FC0")
	BrownColor    = lipgloss.Color("#885D42")
	OrangeColor   = lipgloss.Color("#D8863F")
	BorderColor   = lipgloss.Color("#A77046")
	DangerColor   = lipgloss.Color("#D8863F")
	SuccessColor  = lipgloss.Color("#C7B7A5")
	VanillaColor  = lipgloss.Color("#F2E7D8")
	ChocoColor    = lipgloss.Color("#885D42")

	AppFrameStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(BgBottomColor).
			Padding(1, 2)

	HeroBlueStyle = lipgloss.NewStyle().
			Foreground(BlueColor).
			Bold(true)

	HeroBrownStyle = lipgloss.NewStyle().
			Foreground(BrownColor).
			Bold(true)

	HeroOrangeStyle = lipgloss.NewStyle().
			Foreground(OrangeColor).
			Bold(true)

	SubTitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	ListTitleStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true)

	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(OrangeColor).
				Bold(true)

	CardStyle = lipgloss.NewStyle().
			Background(CardBgColor).
			Border(lipgloss.NormalBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2)

	ActiveCardStyle = lipgloss.NewStyle().
			Background(CardBgColor).
			Border(lipgloss.NormalBorder()).
			BorderForeground(OrangeColor).
			Padding(1, 2)

	PillStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(OrangeColor).
			Padding(0, 2).
			Bold(true)

	MetaStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	InstalledStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	MissingStyle = lipgloss.NewStyle().
			Foreground(DangerColor).
			Bold(true)

	ProjectStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true)

	ActionStyle = lipgloss.NewStyle().
			Foreground(OrangeColor).
			Bold(true)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(OrangeColor).
				Bold(true)

	ItemBodyStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	ItemMetaStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	SelectedRailStyle = lipgloss.NewStyle().
				Foreground(OrangeColor).
				Bold(true)

	HelpPillStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			PaddingRight(2)

	AlertStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DangerColor).
			Foreground(TextColor).
			Background(CardBgColor).
			Padding(1, 2)

	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	InputShellStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(CardBgColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(OrangeColor).
			Padding(0, 1)

	LogoVanillaStyle = lipgloss.NewStyle().
				Foreground(VanillaColor).
				Bold(true)

	LogoChocoStyle = lipgloss.NewStyle().
			Foreground(ChocoColor).
			Bold(true)

	LogoConeStyle = lipgloss.NewStyle().
			Foreground(OrangeColor).
			Bold(true)

	BrandBigStyle = lipgloss.NewStyle().
			Foreground(BrownColor).
			Bold(true)

	BrandMidStyle = lipgloss.NewStyle().
			Foreground(BrownColor).
			Bold(true)

	BrandSmallStyle = lipgloss.NewStyle().
			Foreground(BrownColor).
			Bold(true)
)

func GetCLIColor(id string) lipgloss.TerminalColor {
	return TextColor
}

func RenderHero(state, selectedCLI string) string {
	line1 := HeroBlueStyle.Render("I TESTED")
	line2 := HeroBrownStyle.Render("ESPRESSO")
	line3 := HeroOrangeStyle.Render("BEANS")

	switch state {
	case "list":
		line1 = HeroBlueStyle.Render("PICK AN")
		line2 = HeroBrownStyle.Render("LLM")
		line3 = HeroOrangeStyle.Render("LAUNCHER")
	case "dir_pick":
		line1 = HeroBlueStyle.Render("I PICKED")
		line2 = HeroBrownStyle.Render("A PROJECT")
		line3 = HeroOrangeStyle.Render("READY")
	case "create_dir":
		line1 = HeroBlueStyle.Render("I MADE")
		line2 = HeroBrownStyle.Render("A NEW")
		line3 = HeroOrangeStyle.Render("PROJECT")
	}

	subtitle := "Choose one model, choose one folder under ~/private, and launch into a split iTerm2 tab."
	if selectedCLI != "" && state != "list" {
		subtitle = selectedCLI + " will launch in the selected project directory."
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Background(BgBottomColor).Render(line1),
		lipgloss.NewStyle().Background(BgBottomColor).Render(line2),
		lipgloss.NewStyle().Background(BgBottomColor).Render(line3),
		"",
		lipgloss.NewStyle().Background(BgBottomColor).Render(SubTitleStyle.Render(subtitle)),
	)
}

func RenderLogo() string {
	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		LogoVanillaStyle.Render("o"),
		" ",
		LogoChocoStyle.Render("o"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		LogoConeStyle.Render(" v"),
	)
}

func RenderAmbientDot(glyph string, variant int) string {
	colors := []lipgloss.TerminalColor{VanillaColor, OrangeColor, BrownColor, MutedColor}
	return lipgloss.NewStyle().
		Foreground(colors[variant%len(colors)]).
		Render(glyph)
}

func RenderHelp(parts ...string) string {
	rendered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		rendered = append(rendered, HelpPillStyle.Render(part))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, rendered...)
}

func RenderPanel(content string, active bool, width, height int) string {
	style := CardStyle.Copy()
	if active {
		style = ActiveCardStyle.Copy()
	}
	if width > 0 {
		style = style.Width(width)
	}
	if height > 0 {
		style = style.Height(height)
	}
	return style.Render(content)
}
