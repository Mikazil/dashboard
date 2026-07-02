package ticker

import (
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"

	"dashboard/internal/theme"
)

type RSSFeedConfig struct {
	Name string
	URL  string
}

type Widget struct {
	feeds       []RSSFeedConfig
	headlines   []string
	headlineIdx int
	scrollPos   int
	mu          sync.RWMutex
	updateInt   time.Duration
	scrollSpeed time.Duration
	lastFetch   time.Time
	err         error
	done        chan struct{}
}

func New(feeds []RSSFeedConfig, updateInterval, scrollSpeed time.Duration) *Widget {
	w := &Widget{
		feeds:       feeds,
		updateInt:   updateInterval,
		scrollSpeed: scrollSpeed,
		done:        make(chan struct{}),
	}
	return w
}

func (w *Widget) Init() {
	w.fetchAll()
	go w.scrollLoop()
}

func (w *Widget) fetchAll() {
	var all []string
	parser := gofeed.NewParser()

	for _, feed := range w.feeds {
		f, err := parser.ParseURL(feed.URL)
		if err != nil {
			w.err = err
			continue
		}
		for _, item := range f.Items {
			if item.Title != "" {
				all = append(all, fmt.Sprintf("▸ %s", item.Title))
			}
		}
	}

	w.mu.Lock()
	if len(all) > 0 {
		w.headlines = all
		w.err = nil
	}
	w.lastFetch = time.Now()
	w.mu.Unlock()
}

func (w *Widget) scrollLoop() {
	ticker := time.NewTicker(w.scrollSpeed)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.mu.Lock()
			if len(w.headlines) > 0 {
				w.scrollPos++
				text := w.headlines[w.headlineIdx]
				if w.scrollPos >= len(text) {
					w.headlineIdx = (w.headlineIdx + 1) % len(w.headlines)
					w.scrollPos = 0
				}
			}
			w.mu.Unlock()
		case <-w.done:
			return
		}
	}
}

func (w *Widget) Stop() {
	close(w.done)
}

func (w *Widget) Update() {
	if time.Since(w.lastFetch) >= w.updateInt {
		go w.fetchAll()
	}
}

func (w *Widget) View(width int) string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.err != nil && len(w.headlines) == 0 {
		return theme.Error.Render(" ⚠ RSS error ")
	}

	if len(w.headlines) == 0 {
		return theme.DimText.Render(" Loading... ")
	}

	text := w.headlines[w.headlineIdx]
	sep := "  ◆  "
	track := text + sep + text

	pos := w.scrollPos % (len(text) + len(sep))
	visible := track[pos:]
	if len(visible) > width {
		visible = visible[:width]
	}
	visible = fmt.Sprintf("%-*s", width, visible)

	labelStyle := theme.Title.Copy().
		Background(theme.DimBg).
		Foreground(theme.Primary)

	return lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render(" TECH "),
		theme.Base.Render(" "+visible+" "),
	)
}
