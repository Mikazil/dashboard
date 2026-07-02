package fetcher

import (
	"io"
	"math"
	"net/http"
	"time"
)

type Fetcher struct {
	client  *http.Client
	MaxRetries int
	BaseDelay  time.Duration
}

func New(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		MaxRetries: 3,
		BaseDelay:  time.Second,
	}
}

func (f *Fetcher) Fetch(url string) ([]byte, error) {
	var lastErr error
	for i := range f.MaxRetries {
		resp, err := f.client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				return body, nil
			}
			lastErr = err
		} else {
			lastErr = err
		}

		if i < f.MaxRetries-1 {
			delay := time.Duration(math.Pow(2, float64(i))) * f.BaseDelay
			time.Sleep(delay)
		}
	}
	return nil, lastErr
}
