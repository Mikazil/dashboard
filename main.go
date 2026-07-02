package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"dashboard/internal/config"
	"dashboard/internal/layout"
	"dashboard/internal/theme"
	"dashboard/internal/widgets/calendar"
	"dashboard/internal/widgets/clock"
	"dashboard/internal/widgets/currency"
	"dashboard/internal/widgets/ticker"
	"dashboard/internal/widgets/weather"
)

type tickMsg struct{}

type model struct {
	clock    *clock.Widget
	weather  *weather.Widget
	currency *currency.Widget
	calendar *calendar.Widget
	ticker   *ticker.Widget
	cfg      config.Config
	width    int
	height   int
}

func initialModel() model {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Printf("config: %v, using defaults", err)
	}

	weatherInt, _ := time.ParseDuration(cfg.Weather.Interval)
	if weatherInt == 0 {
		weatherInt = 5 * time.Minute
	}
	currInt, _ := time.ParseDuration(cfg.Currency.Interval)
	if currInt == 0 {
		currInt = 5 * time.Minute
	}
	rssInt, _ := time.ParseDuration(cfg.RSS.UpdateInterval)
	if rssInt == 0 {
		rssInt = time.Hour
	}
	rssSpeed, _ := time.ParseDuration(cfg.RSS.ScrollSpeed)
	if rssSpeed == 0 {
		rssSpeed = 50 * time.Millisecond
	}

	var rssFeeds []ticker.RSSFeedConfig
	for _, f := range cfg.RSS.Feeds {
		rssFeeds = append(rssFeeds, ticker.RSSFeedConfig{
			Name: f.Name,
			URL:  f.URL,
		})
	}
	if len(rssFeeds) == 0 {
		rssFeeds = []ticker.RSSFeedConfig{
			{Name: "Habr", URL: "https://habr.com/ru/rss/hubs/all/"},
		}
	}

	m := model{
		clock:    clock.New(),
		weather:  weather.New(cfg.Weather.City, cfg.Weather.APIKey, weatherInt),
		currency: currency.New(cfg.Currency.Codes, cfg.Currency.Crypto, currInt),
		calendar: calendar.New(),
		ticker:   ticker.New(rssFeeds, rssInt, rssSpeed),
		cfg:      cfg,
		width:    80,
		height:   24,
	}

	m.clock.Init()
	m.ticker.Init()

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.clock.Stop()
			m.ticker.Stop()
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		m.clock.Update()
		m.weather.Update()
		m.currency.Update()
		m.calendar.Update()
		m.ticker.Update()

		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}

	return m, nil
}

func (m model) View() string {
	grid := layout.New([][]layout.Cell{
		{
			{
				Title: "Time",
				View: func(w, h int) string {
					return m.clock.View(w)
				},
			},
			{
				Title: "Weather",
				View: func(w, h int) string {
					return m.weather.View(w)
				},
			},
		},
		{
			{
				Title: "Rates",
				View: func(w, h int) string {
					return m.currency.View(w)
				},
			},
			{
				Title: "Calendar",
				View: func(w, h int) string {
					return m.calendar.View(w)
				},
			},
		},
	})

	topSection := grid.View(m.width, m.height-3)

	bottomStyle := theme.Box.
		Copy().
		Width(m.width - 2).
		UnsetHeight()

	bottom := bottomStyle.Render(m.ticker.View(m.width - 6))

	return theme.Base.Render(topSection + "\n" + bottom)
}

func main() {
	f, err := tea.LogToFile("dashboard.log", "dashboard")
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: %v\n", err)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
