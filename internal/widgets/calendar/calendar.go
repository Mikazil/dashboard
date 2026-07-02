package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"dashboard/internal/theme"
)

type Widget struct {
	now time.Time
}

func New() *Widget {
	return &Widget{}
}

func (w *Widget) Init() {}

func (w *Widget) Update() {
	w.now = time.Now()
}

func (w *Widget) View(width, height int) string {
	now := w.now
	year, month, _ := now.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	weekday := firstDay.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	offset := int(weekday) - 1
	if offset < 0 {
		offset = 6
	}

	var sb strings.Builder
	sb.WriteString(theme.DimText.Render(fmt.Sprintf(" %s %d ", month.String()[:3], year)) + "\n")
	sb.WriteString(theme.DimText.Render(" Mo Tu We Th Fr Sa Su") + "\n")

	var lines []string
	var line strings.Builder

	for range offset {
		line.WriteString("   ")
	}

	for day := 1; day <= lastDay.Day(); day++ {
		d := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		isToday := d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day()

		var cell string
		if isToday {
			highlight := theme.Base.Copy().
				Background(theme.Secondary).
				Foreground(theme.Bg)
			cell = highlight.Render(fmt.Sprintf("%2d", day))
		} else {
			cell = theme.DimText.Render(fmt.Sprintf("%2d", day))
		}

		if line.Len() > 0 {
			line.WriteString("  ")
		}
		line.WriteString(cell)

		dow := d.Weekday()
		if dow == time.Sunday && day != lastDay.Day() {
			lines = append(lines, line.String())
			line.Reset()
		}
	}
	if line.Len() > 0 {
		lines = append(lines, line.String())
	}

	for _, l := range lines {
		sb.WriteString(l + "\n")
	}

	content := lipgloss.NewStyle().Width(width - 2).Render(sb.String())
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
