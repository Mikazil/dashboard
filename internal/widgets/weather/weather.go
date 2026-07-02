package weather

import (
	"fmt"
	"strings"
	"time"

	"dashboard/internal/fetcher"
	"dashboard/internal/theme"
)

type WeatherData struct {
	Temperature string
	Description string
	Humidity    string
	Wind        string
	Pressure    string
	Icon        string
}

type Widget struct {
	data      *WeatherData
	ftch      *fetcher.Fetcher
	city      string
	interval  time.Duration
	lastFetch time.Time
	err       error
}

func New(city string, interval time.Duration) *Widget {
	return &Widget{
		ftch:     fetcher.New(15 * time.Second),
		city:     city,
		interval: interval,
	}
}

func (w *Widget) Init() {}

func (w *Widget) Update() {
	if time.Since(w.lastFetch) < w.interval {
		return
	}
	w.fetch()
}

func (w *Widget) fetch() {
	body, err := w.ftch.Fetch(fmt.Sprintf("https://wttr.in/%s?format=%%t+%%h+%%w+%%P+%%C", w.city))
	if err != nil {
		w.err = err
		return
	}

	line := strings.TrimSpace(string(body))
	parts := strings.Split(line, " ")
	if len(parts) < 5 {
		w.err = fmt.Errorf("unexpected format: %s", line)
		return
	}

	temp := parts[0]
	hum := strings.TrimRight(parts[1], "%")
	wind := parts[2]
	pres := parts[3]
	desc := strings.Join(parts[4:], " ")

	w.data = &WeatherData{
		Temperature: temp,
		Description: desc,
		Humidity:    hum + "%",
		Wind:        wind,
		Pressure:    pres,
		Icon:        getASCIIIcon(desc),
	}
	w.err = nil
	w.lastFetch = time.Now()
}

func (w *Widget) View(width int) string {
	if w.err != nil {
		return theme.Error.Render(" ⚠ Weather error ")
	}

	if w.data == nil {
		return theme.DimText.Render(" Loading... ")
	}

	return fmt.Sprintf("%s\n%s\n %s\n %s\n\nHumidity %s\nWind     %s\nPressure %s",
		theme.DimText.Render(" Weather "),
		theme.Base.Render(w.data.Icon),
		theme.Base.Render(w.data.Temperature),
		theme.DimText.Render(w.data.Description),
		theme.Base.Render(w.data.Humidity),
		theme.Base.Render(w.data.Wind),
		theme.Base.Render(w.data.Pressure),
	)
}

func getASCIIIcon(desc string) string {
	desc = strings.ToLower(desc)
	switch {
	case strings.Contains(desc, "sun"), strings.Contains(desc, "clear"):
		return "   \\   /\n    \\ /\n  .--.--.\n /  _    \\\n|  / \\   |\n|  \\_/   |\n \\      /\n  `----`"
	case strings.Contains(desc, "cloud"), strings.Contains(desc, "overcast"), strings.Contains(desc, "mist"):
		return "    .--.\n .-(    ).\n(___(__)__)\n           "
	case strings.Contains(desc, "rain"), strings.Contains(desc, "drizzle"):
		return "    _  _\n   /  _/\n  / _/\n /_/\n\\   \\\n \\   \\\n  \\\n   \\"
	case strings.Contains(desc, "snow"), strings.Contains(desc, "blizzard"):
		return "   *  *\n  * ** *\n  ** **\n   * *\n  * ***\n  ** **"
	default:
		return "    .--.\n .-(    ).\n(___(__)__)\n           "
	}
}
