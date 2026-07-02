package clock

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"dashboard/internal/theme"
)

type Widget struct {
	time time.Time
}

func New() *Widget {
	return &Widget{}
}

func (w *Widget) Init() {
	w.time = time.Now()
}

func (w *Widget) Stop() {}

func (w *Widget) Update() {
	w.time = time.Now()
}

func (w *Widget) View(width, height int) string {
	now := w.time
	loc, _ := time.LoadLocation("Europe/Moscow")
	if loc != nil {
		now = now.In(loc)
	}

	timeStr := now.Format("15:04:05")
	dateStr := now.Format("Monday")
	dateStr2 := now.Format("02 Jan 2006")

	clock := bigClock(timeStr, width)

	pad1 := (width - lipgloss.Width(dateStr)) / 2
	if pad1 < 0 {
		pad1 = 0
	}
	date1 := strings.Repeat(" ", pad1) + theme.DimText.Render(dateStr)

	pad2 := (width - lipgloss.Width(dateStr2)) / 2
	if pad2 < 0 {
		pad2 = 0
	}
	date2 := strings.Repeat(" ", pad2) + theme.DimText.Render(dateStr2)

	content := lipgloss.JoinVertical(lipgloss.Center, clock, date1, date2)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func bigClock(t string, width int) string {
	digits := [][]string{
		{" ██ ", "█  █", "█  █", "█  █", " ██ "},
		{"  █ ", "  █ ", "  █ ", "  █ ", "  █ "},
		{" ██ ", "   █", " ██ ", "█   ", " ██ "},
		{" ██ ", "   █", " ██ ", "   █", " ██ "},
		{"█  █", "█  █", " ██ ", "   █", "   █"},
		{" ██ ", "█   ", " ██ ", "   █", " ██ "},
		{" ██ ", "█   ", " ██ ", "█  █", " ██ "},
		{"███ ", "   █", "  █ ", " █  ", " █  "},
		{" ██ ", "█  █", " ██ ", "█  █", " ██ "},
		{" ██ ", "█  █", " ██ ", "   █", " ██ "},
		{"    ", " ██ ", "    ", " ██ ", "    "},
	}

	var lines [5]string

	for _, ch := range t {
		var idx int
		if ch == ':' {
			idx = 10
		} else if ch >= '0' && ch <= '9' {
			idx = int(ch - '0')
		} else {
			continue
		}
		for i := range 5 {
			lines[i] += digits[idx][i]
		}
	}

	clockWidth := lipgloss.Width(lines[0])
	padLeft := (width - clockWidth) / 2
	if padLeft < 0 {
		padLeft = 0
	}

	result := make([]string, 5)
	for i := range lines {
		styled := theme.BigText.Copy().Foreground(theme.Primary)
		result[i] = styled.Render(strings.Repeat(" ", padLeft) + lines[i])
	}

	return strings.Join(result, "\n")
}
