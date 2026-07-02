package clock

import (
	"time"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"dashboard/internal/theme"
)

type Widget struct {
	time   time.Time
	ticker *time.Ticker
	done   chan struct{}
}

func New() *Widget {
	return &Widget{
		ticker: time.NewTicker(time.Second),
		done:   make(chan struct{}),
	}
}

func (w *Widget) Init() {
	go func() {
		for {
			select {
			case t := <-w.ticker.C:
				w.time = t
			case <-w.done:
				return
			}
		}
	}()
}

func (w *Widget) Stop() {
	w.ticker.Stop()
	close(w.done)
}

func (w *Widget) Update() {}

func (w *Widget) View(width int) string {
	now := time.Now()
	loc, _ := time.LoadLocation("Europe/Moscow")
	if loc != nil {
		now = now.In(loc)
	}

	timeStr := now.Format("15:04:05")
	dateStr := now.Format("Monday")
	dateStr2 := now.Format("02 Jan 2006")

	timeW := lipgloss.Width(timeStr)
	pad := (width - timeW) / 2
	if pad < 0 {
		pad = 0
	}
	timeLine := strings.Repeat(" ", pad) + theme.BigText.Render(timeStr)

	date1W := lipgloss.Width(dateStr)
	pad1 := (width - date1W) / 2
	if pad1 < 0 {
		pad1 = 0
	}
	dateLine1 := strings.Repeat(" ", pad1) + theme.DimText.Render(dateStr)

	date2W := lipgloss.Width(dateStr2)
	pad2 := (width - date2W) / 2
	if pad2 < 0 {
		pad2 = 0
	}
	dateLine2 := strings.Repeat(" ", pad2) + theme.DimText.Render(dateStr2)

	clockArt := bigClock(timeStr)
	if clockArt != "" {
		return lipgloss.JoinVertical(lipgloss.Center,
			clockArt,
			dateLine1,
			dateLine2,
		)
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		timeLine,
		dateLine1,
		dateLine2,
	)
}

func bigClock(t string) string {
	digits := [][]string{
		{" ██ ", "█  █", "█  █", "█  █", " ██ "}, // 0
		{"  █ ", " ██ ", "  █ ", "  █ ", " ███"}, // 1
		{" ██ ", "█  █", "  █ ", " █  ", "████"}, // 2
		{" ██ ", "█  █", "  █ ", "█  █", " ██ "}, // 3
		{"█  █", "█  █", "████", "   █", "   █"}, // 4
		{"████", "█   ", "███ ", "   █", "███ "}, // 5
		{" ██ ", "█   ", "███ ", "█  █", " ██ "}, // 6
		{"████", "   █", "  █ ", " █  ", " █  "}, // 7
		{" ██ ", "█  █", " ██ ", "█  █", " ██ "}, // 8
		{" ██ ", "█  █", " ██ ", "   █", " ██ "}, // 9
		{"    ", " ░░ ", "    ", " ░░ ", "    "}, // :
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

	result := make([]string, 5)
	for i := range lines {
		style := theme.BigText.Copy().Foreground(theme.Primary)
		result[i] = style.Render(lines[i])
	}

	return strings.Join(result, "\n")
}
