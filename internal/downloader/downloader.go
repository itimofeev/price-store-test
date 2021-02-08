package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func New() *Downloader {
	return &Downloader{
		client: &http.Client{Transport: http.DefaultTransport},
	}
}

type Downloader struct {
	client *http.Client
}

// GetCSV downloads csv by passed URL
// client of this method is responsible for closing returned io.ReadCloser
func (d *Downloader) GetCSV(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}
	resp, err := d.client.Do(req) //nolint:bodyclose
	if err != nil {
		return nil, fmt.Errorf("get url failed: %w", err)
	}

	return resp.Body, nil
}
