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

func (w *Widget) View(width int) string {
	now := w.time
	loc, _ := time.LoadLocation("Europe/Moscow")
	if loc != nil {
		now = now.In(loc)
	}

	timeStr := now.Format("15:04:05")
	dateStr := now.Format("Monday")
	dateStr2 := now.Format("02 Jan 2006")

	pad1 := (width - lipgloss.Width(dateStr)) / 2
	if pad1 < 0 {
		pad1 = 0
	}
	dateLine1 := strings.Repeat(" ", pad1) + theme.DimText.Render(dateStr)

	pad2 := (width - lipgloss.Width(dateStr2)) / 2
	if pad2 < 0 {
		pad2 = 0
	}
	dateLine2 := strings.Repeat(" ", pad2) + theme.DimText.Render(dateStr2)

	return lipgloss.JoinVertical(lipgloss.Center,
		bigClock(timeStr),
		dateLine1,
		dateLine2,
	)
}

func bigClock(t string) string {
	digits := [][]string{
		{" ██ ", "█  █", "█  █", "█  █", " ██ "},
		{"  █ ", " ██ ", "  █ ", "  █ ", " ███"},
		{" ██ ", "█  █", "  █ ", " █  ", "████"},
		{" ██ ", "█  █", "  █ ", "█  █", " ██ "},
		{"█  █", "█  █", "████", "   █", "   █"},
		{"████", "█   ", "███ ", "   █", "███ "},
		{" ██ ", "█   ", "███ ", "█  █", " ██ "},
		{"████", "   █", "  █ ", " █  ", " █  "},
		{" ██ ", "█  █", " ██ ", "█  █", " ██ "},
		{" ██ ", "█  █", " ██ ", "   █", " ██ "},
		{"    ", " ░░ ", "    ", " ░░ ", "    "},
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
			lines[i] += digits[idx][i] + " "
		}
	}

	result := make([]string, 5)
	for i := range lines {
		style := theme.BigText.Copy().Foreground(theme.Primary)
		result[i] = style.Render(lines[i])
	}

	return strings.Join(result, "\n")
}
