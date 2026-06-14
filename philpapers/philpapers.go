// Package philpapers is the library behind the philpapers command line:
// the HTTP client, request shaping, and the typed data models for PhilPapers
// philosophy papers fetched via the OAI-PMH endpoint.
package philpapers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const DefaultUserAgent = "Mozilla/5.0 (compatible; philpapers-cli/0.1; +https://github.com/tamnd/philpapers-cli)"

// Config holds constructor parameters.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://philpapers.org",
		UserAgent: DefaultUserAgent,
		Rate:      1 * time.Second,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to the PhilPapers OAI-PMH endpoint.
type Client struct {
	cfg        Config
	httpClient *http.Client
	mu         sync.Mutex
	last       time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) get(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		b, retry, err := c.do(ctx, rawURL)
		if err == nil {
			return b, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get: %w", lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}

// Recent fetches recently added PhilPapers entries via OAI-PMH ListRecords.
func (c *Client) Recent(ctx context.Context, limit int) ([]Paper, error) {
	if limit <= 0 {
		limit = 20
	}
	apiURL := c.cfg.BaseURL + "/oai.pl?verb=ListRecords&metadataPrefix=oai_dc"
	raw, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}
	return parseOAI(raw, limit), nil
}

func parseOAI(raw []byte, limit int) []Paper {
	var doc oaiPMH
	if err := xml.Unmarshal(raw, &doc); err != nil || doc.ListRecords == nil {
		return nil
	}
	var out []Paper
	for i, rec := range doc.ListRecords.Records {
		if limit > 0 && i >= limit {
			break
		}
		dc := rec.Metadata.DC
		title := ""
		if len(dc.Titles) > 0 {
			title = dc.Titles[0]
		}
		author := ""
		if len(dc.Creators) > 0 {
			author = dc.Creators[0]
		}
		date := ""
		if len(dc.Dates) > 0 {
			date = dc.Dates[0]
		}
		subject := ""
		if len(dc.Subjects) > 0 {
			subject = dc.Subjects[0]
		}
		url := ""
		for _, r := range dc.Relations {
			if strings.HasPrefix(r, "http") {
				url = r
				break
			}
		}
		// Extract ID from identifier: "oai:philpapers.org/rec/SMIT-123" -> "SMIT-123"
		id := rec.Header.Identifier
		if idx := strings.LastIndex(id, "/"); idx >= 0 {
			id = id[idx+1:]
		}
		out = append(out, Paper{
			Rank:    i + 1,
			ID:      id,
			Title:   title,
			Author:  author,
			Date:    date,
			Subject: subject,
			URL:     url,
		})
	}
	return out
}
