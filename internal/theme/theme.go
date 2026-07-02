package theme

import "github.com/charmbracelet/lipgloss"

var (
	Primary   = lipgloss.Color("#FFB000")
	Secondary = lipgloss.Color("#E89900")
	Dim       = lipgloss.Color("#886600")
	DimBg     = lipgloss.Color("#221100")
	Bg        = lipgloss.Color("#000000")

	Base = lipgloss.NewStyle().
		Background(Bg).
		Foreground(Primary)

	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Dim).
		Background(Bg).
		Foreground(Primary).
		Padding(0, 1)

	Header = lipgloss.NewStyle().
		Background(DimBg).
		Foreground(Primary).
		Bold(true)

	Ticker  = lipgloss.NewStyle().Background(Bg).Foreground(Secondary)
	Error   = lipgloss.NewStyle().Background(Bg).Foreground(lipgloss.Color("#CC4400"))
	Success = lipgloss.NewStyle().Background(Bg).Foreground(Primary)

	Title = lipgloss.NewStyle().
		Background(DimBg).
		Foreground(Primary).
		Bold(true).
		Padding(0, 1, 0, 1)

	BigText = lipgloss.NewStyle().
		Background(Bg).
		Foreground(Primary).
		Bold(true)

	DimText = lipgloss.NewStyle().
		Background(Bg).
		Foreground(Dim)
)
