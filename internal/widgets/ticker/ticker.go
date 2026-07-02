package ticker

import (
	"sync"
	"time"
	"unicode/utf8"

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
	track        []rune
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
	if len(w.headlines) == 0 {
		w.track = nil
		return
	}
	var result string
	for _, h := range w.headlines {
		result += " > " + h + " | "
	}
	w.track = []rune(result)
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
			w.mu.Lock()
			w.buildTrack()
			w.mu.Unlock()
		}()
	}
}

func (w *Widget) View(width int) string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.err != nil && len(w.headlines) == 0 {
		return theme.Error.Render(" ⚠ RSS error ")
	}

	if len(w.headlines) == 0 || len(w.track) == 0 {
		return theme.DimText.Render(" Loading... ")
	}

	trackLen := len(w.track)
	pos := w.scrollPos % trackLen

	var visible []rune
	if pos+width <= trackLen {
		visible = w.track[pos : pos+width]
	} else {
		first := w.track[pos:]
		second := w.track[:width-len(first)]
		visible = append(first, second...)
	}

	text := string(visible)
	if utf8.RuneCountInString(text) < width {
		need := width - utf8.RuneCountInString(text)
		if need > trackLen {
			need = trackLen
		}
		text += string(w.track[:need])
	}

	labelStyle := theme.Title.Copy().
		Background(theme.DimBg).
		Foreground(theme.Primary)

	return lipgloss.JoinHorizontal(lipgloss.Center,
		labelStyle.Render(" TECH "),
		theme.Base.Render(" "+text+" "),
	)
}
