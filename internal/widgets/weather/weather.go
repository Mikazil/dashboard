package weather

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

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
	apiKey    string
	interval  time.Duration
	lastFetch time.Time
	err       error
}

func New(city, apiKey string, interval time.Duration) *Widget {
	return &Widget{
		ftch:     fetcher.New(15 * time.Second),
		city:     city,
		apiKey:   apiKey,
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
	if w.apiKey == "" {
		w.err = fmt.Errorf("no api_key in config")
		return
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", w.city, w.apiKey)
	body, err := w.ftch.Fetch(url)
	if err != nil {
		w.err = err
		return
	}

	var resp struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity float64 `json:"humidity"`
			Pressure float64 `json:"pressure"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		w.err = err
		return
	}

	if len(resp.Weather) == 0 {
		w.err = fmt.Errorf("no weather data")
		return
	}

	desc := resp.Weather[0].Description
	pres := int(math.Round(resp.Main.Pressure * 0.75006))

	w.data = &WeatherData{
		Temperature: fmt.Sprintf("%.0f°C", resp.Main.Temp),
		Description: strings.Title(desc),
		Humidity:    fmt.Sprintf("%.0f%%", resp.Main.Humidity),
		Wind:        fmt.Sprintf("%.0fm/s", resp.Wind.Speed),
		Pressure:    fmt.Sprintf("%dmmHg", pres),
		Icon:        getASCIIIcon(desc + " " + resp.Weather[0].Main),
	}
	w.err = nil
	w.lastFetch = time.Now()
}

func (w *Widget) View(width, height int) string {
	var content string

	if w.err != nil {
		errMsg := w.err.Error()
		if len(errMsg) > 30 {
			errMsg = errMsg[:30]
		}
		content = theme.Error.Render(" ! " + errMsg + " ")
	} else if w.data == nil {
		content = theme.DimText.Render(" Loading... ")
	} else {
		content = fmt.Sprintf("%s\n%s\n %s\n %s\n\nHumidity %s\nWind     %s\nPressure %s",
			theme.DimText.Render(" Weather "),
			theme.Base.Render(w.data.Icon),
			theme.Base.Render(w.data.Temperature),
			theme.DimText.Render(w.data.Description),
			theme.Base.Render(w.data.Humidity),
			theme.Base.Render(w.data.Wind),
			theme.Base.Render(w.data.Pressure),
		)
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
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
