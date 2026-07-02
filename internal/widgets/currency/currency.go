package currency

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"dashboard/internal/fetcher"
	"dashboard/internal/theme"
)

type Rate struct {
	Code   string
	Name   string
	Value  float64
	Change float64
}

type Widget struct {
	rates    []Rate
	cryptos  []Rate
	fetcher  *fetcher.Fetcher
	codes    []string
	cryptoSymbols []string
	interval time.Duration
	lastFetch time.Time
	err      error
}

func New(codes, cryptoSymbols []string, interval time.Duration) *Widget {
	return &Widget{
		fetcher:  fetcher.New(15 * time.Second),
		codes:    codes,
		cryptoSymbols: cryptoSymbols,
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
	w.fetchFiat()
	w.fetchCrypto()
	w.lastFetch = time.Now()
}

func (w *Widget) fetchFiat() {
	oldRates := w.rates

	resp, err := http.Get("https://www.cbr.ru/scripts/XML_daily.asp")
	if err != nil {
		w.err = err
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.err = err
		return
	}

	type Valute struct {
		CharCode string `xml:"CharCode"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	}

	type ValCurs struct {
		Valutes []Valute `xml:"Valute"`
	}

	var vc ValCurs
	if err := xml.Unmarshal(body, &vc); err != nil {
		w.err = err
		return
	}

	valMap := make(map[string]float64)
	for _, v := range vc.Valutes {
		valStr := strings.Replace(v.Value, ",", ".", 1)
		if val, err := strconv.ParseFloat(valStr, 64); err == nil {
			valMap[v.CharCode] = val
		}
	}

	var newRates []Rate
	for _, code := range w.codes {
		if val, ok := valMap[code]; ok {
			oldVal := 0.0
			for _, r := range oldRates {
				if r.Code == code {
					oldVal = r.Value
					break
				}
			}
			newRates = append(newRates, Rate{
				Code:   code,
				Value:  val,
				Change: val - oldVal,
			})
		}
	}
	if len(newRates) > 0 {
		w.rates = newRates
	}
}

func (w *Widget) fetchCrypto() {
	oldRates := w.cryptos

	ids := strings.Join(w.cryptoSymbols, ",")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", ids)
	body, err := w.fetcher.Fetch(url)
	if err != nil {
		w.err = err
		return
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		w.err = err
		return
	}

	var newCryptos []Rate
	for _, id := range w.cryptoSymbols {
		if prices, ok := result[id]; ok {
			if price, ok := prices["usd"]; ok {
				oldVal := 0.0
				for _, r := range oldRates {
					if r.Code == id {
						oldVal = r.Value
						break
					}
				}
				newCryptos = append(newCryptos, Rate{
					Code:   strings.ToUpper(id[:3]),
					Value:  price,
					Change: price - oldVal,
				})
			}
		}
	}
	if len(newCryptos) > 0 {
		w.cryptos = newCryptos
	}
}

func (w *Widget) View(width int) string {
	var sb strings.Builder

	if w.err != nil {
		sb.WriteString(theme.Error.Render(" ⚠ Rates error "))
	} else {
		sb.WriteString(theme.DimText.Render(" Exchange rates ") + "\n")
		for _, r := range w.rates {
			arrow := " "
			if r.Change > 0 {
				arrow = "↑"
			} else if r.Change < 0 {
				arrow = "↓"
			}
			valStr := fmt.Sprintf("%.2f", r.Value)
			sb.WriteString(fmt.Sprintf("%s  %s %s", r.Code, arrow, valStr) + "\n")
		}

		if len(w.cryptos) > 0 {
			sb.WriteString("")
			for _, c := range w.cryptos {
				arrow := " "
				if c.Change > 0 {
					arrow = "↑"
				} else if c.Change < 0 {
					arrow = "↓"
				}
				sb.WriteString(fmt.Sprintf("%s %s $%.0f", c.Code, arrow, c.Value) + "\n")
			}
		}
	}

	return lipgloss.NewStyle().
		Width(width - 2).
		Render(sb.String())
}
