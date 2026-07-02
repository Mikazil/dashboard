package ticker

import (
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
	feeds        []RSSFeedConfig
	headlines    []string
	track        string
	scrollPos    int
	mu           sync.RWMutex
	updateInt    time.Duration
	scrollSpeed  time.Duration
	lastFetch    time.Time
	err          error
	done         chan struct{}
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
	w.buildTrack()
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
				all = append(all, item.Title)
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

func (w *Widget) buildTrack() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.headlines) == 0 {
		w.track = ""
		return
	}

	var result string
	for _, h := range w.headlines {
		result += " ▸ " + h + "  ◆ "
	}
	result += " "
	w.track = result
}

func (w *Widget) scrollLoop() {
	ticker := time.NewTicker(w.scrollSpeed)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.mu.Lock()
			w.scrollPos++
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
		go func() {
			w.fetchAll()
			w.buildTrack()
		}()
	}
}

func (w *Widget) View(width int) string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.err != nil && len(w.headlines) == 0 {
		return theme.Error.Render(" ⚠ RSS error ")
	}

	if len(w.headlines) == 0 || w.track == "" {
		return theme.DimText.Render(" Loading... ")
	}

	pos := w.scrollPos % len(w.track)
	end := pos + width
	if end > len(w.track) {
		end = len(w.track)
	}
	visible := w.track[pos:end]
	if len(visible) < width {
		visible += w.track[:width-len(visible)]
	}

	labelStyle := theme.Title.Copy().
		Background(theme.DimBg).
		Foreground(theme.Primary)

	return lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render(" TECH "),
		theme.Base.Render(" "+visible+" "),
	)
}
