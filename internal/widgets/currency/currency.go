package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"dashboard/internal/fetcher"
	"dashboard/internal/theme"
)

type Rate struct {
	Code   string
	Value  float64
	Change float64
}

type Widget struct {
	rates         []Rate
	cryptos       []Rate
	ftch          *fetcher.Fetcher
	codes         []string
	cryptoSymbols []string
	cryptoNames   map[string]string
	interval      time.Duration
	lastFetch     time.Time
	err           error
}

func New(codes, cryptoSymbols []string, interval time.Duration) *Widget {
	cryptoNames := map[string]string{
		"bitcoin":     "BTC",
		"ethereum":    "ETH",
		"tether":      "USDT",
		"solana":      "SOL",
		"cardano":     "ADA",
	}
	return &Widget{
		ftch:          fetcher.New(15 * time.Second),
		codes:         codes,
		cryptoSymbols: cryptoSymbols,
		cryptoNames:   cryptoNames,
		interval:      interval,
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

	resp, err := http.Get("https://open.er-api.com/v6/latest/RUB")
	if err != nil {
		w.err = err
		return
	}
	defer resp.Body.Close()

	var result struct {
		Base  string             `json:"base_code"`
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		w.err = err
		return
	}

	var newRates []Rate
	for _, code := range w.codes {
		rate, ok := result.Rates[code]
		if !ok {
			continue
		}
		// API returns how many units per 1 RUB.
		// Convert to RUB per unit (standard display)
		val := 1.0 / rate

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
	if len(newRates) > 0 {
		w.rates = newRates
		w.err = nil
	}
}

func (w *Widget) fetchCrypto() {
	if len(w.cryptoSymbols) == 0 {
		return
	}
	oldCryptos := w.cryptos

	ids := strings.Join(w.cryptoSymbols, ",")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", ids)
	body, err := w.ftch.Fetch(url)
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
				for _, r := range oldCryptos {
					if r.Code == id {
						oldVal = r.Value
						break
					}
				}
				code := w.cryptoNames[id]
				if code == "" {
					code = strings.ToUpper(id)
				}
				newCryptos = append(newCryptos, Rate{
					Code:   code,
					Value:  price,
					Change: price - oldVal,
				})
			}
		}
	}
	if len(newCryptos) > 0 {
		w.cryptos = newCryptos
		w.err = nil
	}
}

func (w *Widget) View(width int) string {
	var sb strings.Builder

	if w.err != nil && len(w.rates) == 0 {
		sb.WriteString(theme.Error.Render(" ⚠ Rates error "))
	}

	if len(w.rates) == 0 {
		sb.WriteString(theme.DimText.Render(" Loading rates... "))
	} else {
		for _, r := range w.rates {
			arrow := " "
			if r.Change > 0 {
				arrow = "↑"
			} else if r.Change < 0 {
				arrow = "↓"
			}
			sb.WriteString(fmt.Sprintf("%s  %s %.2f", r.Code, arrow, r.Value) + "\n")
		}
	}

	if len(w.cryptos) > 0 {
		sb.WriteString("\n")
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

	return lipgloss.NewStyle().Width(width - 2).Render(sb.String())
}
