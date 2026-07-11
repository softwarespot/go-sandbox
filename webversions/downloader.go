package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Option func(c *downloaderConfig)

// WithTimeout sets the timeout for the downloader.
func WithTimeout(timeout time.Duration) Option {
	return func(c *downloaderConfig) {
		c.timeout = timeout
	}
}

// WithUserAgent sets the User-Agent header for the downloader.
func WithUserAgent(userAgent string) Option {
	return func(c *downloaderConfig) {
		if userAgent == "" {
			panic("userAgent cannot be empty")
		}
		c.userAgent = userAgent
	}
}

type downloaderConfig struct {
	timeout   time.Duration
	userAgent string
}

type Downloader struct {
	hc     *http.Client
	config *downloaderConfig
}

// NewDownloader creates a new Downloader with a default HTTP client.
func NewDownloader(opts ...Option) *Downloader {
	cfg := &downloaderConfig{
		timeout:   30 * time.Second,
		userAgent: "WebVersions/1.0",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Downloader{
		hc: &http.Client{
			Timeout: cfg.timeout,
		},
		config: cfg,
	}
}

// Download downloads the content from the provided URL and returns it as a string.
func (d *Downloader) Download(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("User-Agent", d.config.userAgent)
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
