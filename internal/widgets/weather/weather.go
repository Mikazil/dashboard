package weather

import (
	"encoding/json"
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
	data     *WeatherData
	fetcher  *fetcher.Fetcher
	city     string
	interval time.Duration
	lastFetch time.Time
	err      error
}

func New(city string, interval time.Duration) *Widget {
	return &Widget{
		fetcher:  fetcher.New(15 * time.Second),
		city:     city,
		interval: interval,
	}
}

func (w *Widget) Init() {}

func (w *Widget) Update() {
	if time.Since(w.lastFetch) < w.interval {
		return
	}
	w.fetchWeather()
}

func (w *Widget) fetchWeather() {
	body, err := w.fetcher.Fetch(fmt.Sprintf("https://wttr.in/%s?format=j1", w.city))
	if err != nil {
		w.err = err
		return
	}

	var resp struct {
		CurrentCondition []struct {
			TempC       string `json:"temp_C"`
			Humidity    string `json:"humidity"`
			WindSpeed   string `json:"windspeedKmph"`
			Pressure    string `json:"pressureMB"`
			WeatherDesc []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
		} `json:"current_condition"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		w.err = err
		return
	}

	if len(resp.CurrentCondition) == 0 {
		w.err = fmt.Errorf("no weather data")
		return
	}

	cc := resp.CurrentCondition[0]
	desc := ""
	if len(cc.WeatherDesc) > 0 {
		desc = cc.WeatherDesc[0].Value
	}

	w.data = &WeatherData{
		Temperature: cc.TempC + "°C",
		Description: desc,
		Humidity:    cc.Humidity + "%",
		Wind:        cc.WindSpeed + "m/s",
		Pressure:    cc.Pressure + "mmHg",
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
	case strings.Contains(desc, "cloud"), strings.Contains(desc, "overcast"):
		return "    .--.\n .-(    ).\n(___(__)__)\n           "
	case strings.Contains(desc, "rain"), strings.Contains(desc, "drizzle"):
		return "    _  _\n   /  _/\n  / _/\n /_/\n\\   \\\n \\   \\\n  \\\n   \\"
	case strings.Contains(desc, "snow"):
		return "   *  *\n  * ** *\n  ** **\n   * *\n  * ***\n  ** **"
	default:
		return "    .--.\n .-(    ).\n(___(__)__)\n           "
	}
}
