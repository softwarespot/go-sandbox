package downloader

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Option is a function that configures the downloader.
type Option func(c *config)

// WithTimeout sets the timeout for the downloader.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
	}
}

// WithUserAgent sets the User-Agent header for the downloader.
func WithUserAgent(userAgent string) Option {
	return func(c *config) {
		if userAgent == "" {
			panic("userAgent must not be empty")
		}
		c.userAgent = userAgent
	}
}

type config struct {
	timeout   time.Duration
	userAgent string
}

type Downloader struct {
	hc     *http.Client
	config *config
}

// New creates a downloader.
func New(opts ...Option) *Downloader {
	cfg := &config{
		timeout:   30 * time.Second,
		userAgent: "",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	d := &Downloader{
		hc: &http.Client{
			Timeout: cfg.timeout,
		},
		config: cfg,
	}
	return d
}

// Download downloads the content from the provided URL and returns it as a string.
func (d *Downloader) Download(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w", err)
	}
	if d.config.userAgent != "" {
		req.Header.Set("User-Agent", d.config.userAgent)
	}
	resp, err := d.hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("http.Do: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll: %w", err)
	}
	return string(body), nil
}
