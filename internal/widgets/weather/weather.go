package weather

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"dashboard/internal/fetcher"
	"dashboard/internal/theme"
)

var coords = map[string][2]float64{
	"Moscow":    {55.75, 37.62},
	"Saint Petersburg": {59.93, 30.33},
	"Novosibirsk": {55.04, 82.93},
	"Yekaterinburg": {56.84, 60.65},
	"Kazan":     {55.79, 49.12},
	"London":    {51.51, -0.13},
	"New York":  {40.71, -74.01},
	"Tokyo":     {35.68, 139.69},
}

func latLon(city string) (float64, float64) {
	if c, ok := coords[city]; ok {
		return c[0], c[1]
	}
	return 55.75, 37.62
}

func wmoDesc(code int) string {
	switch {
	case code == 0:
		return "Clear"
	case code <= 3:
		return "Cloudy"
	case code <= 48:
		return "Fog"
	case code <= 55:
		return "Drizzle"
	case code <= 65:
		return "Rain"
	case code <= 75:
		return "Snow"
	case code <= 82:
		return "Showers"
	case code >= 95:
		return "Storm"
	default:
		return "Cloudy"
	}
}

func wmoIcon(code int) string {
	switch {
	case code == 0:
		return "   \\   /\n    \\ /\n  .--.--.\n /  _    \\\n|  / \\   |\n|  \\_/   |\n \\      /\n  `----`"
	case code <= 3:
		return "    .--.\n .-(    ).\n(___(__)__)\n           "
	case code <= 48:
		return "    .--.\n .-(    ).\n(___(__)__)\n  // ///\n // // //"
	case code <= 55:
		return "    _  _\n   /  _/\n  / _/\n /_/\n\\   \\\n \\   \\\n  \\\n   \\"
	case code <= 65:
		return "    _  _\n   /  _/\n  / _/\n /_/\n\\   \\\n \\   \\\n  \\\n   \\"
	case code <= 75:
		return "   *  *\n  * ** *\n  ** **\n   * *\n  * ***\n  ** **"
	case code <= 82:
		return "    .--.\n .-(    ).\n(___(__)__)\n  // ///\n // // //"
	default:
		return "  ⚡\n .--.\n.(    ).\n(___(__)__)\n  // ///\n // // //"
	}
}

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
	lat, lon  float64
	city      string
	interval  time.Duration
	lastFetch time.Time
	err       error
}

func New(city string, interval time.Duration) *Widget {
	lat, lon := latLon(city)
	return &Widget{
		ftch:     fetcher.New(15 * time.Second),
		city:     city,
		lat:      lat,
		lon:      lon,
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
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m,pressure_msl&timezone=auto",
		w.lat, w.lon,
	)
	body, err := w.ftch.Fetch(url)
	if err != nil {
		w.err = err
		return
	}

	var resp struct {
		Current struct {
			Temp      float64 `json:"temperature_2m"`
			Humidity  float64 `json:"relative_humidity_2m"`
			Code      int     `json:"weather_code"`
			Wind      float64 `json:"wind_speed_10m"`
			Pressure  float64 `json:"pressure_msl"`
		} `json:"current"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		w.err = err
		return
	}

	desc := wmoDesc(resp.Current.Code)
	w.data = &WeatherData{
		Temperature: fmt.Sprintf("%.0f°C", resp.Current.Temp),
		Description: desc,
		Humidity:    fmt.Sprintf("%.0f%%", resp.Current.Humidity),
		Wind:        fmt.Sprintf("%.0fm/s", resp.Current.Wind),
		Pressure:    fmt.Sprintf("%.0fmmHg", math.Round(resp.Current.Pressure*0.75006)),
		Icon:        wmoIcon(resp.Current.Code),
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
