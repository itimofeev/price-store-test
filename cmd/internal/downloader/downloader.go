package downloader

import (
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
func (d *Downloader) GetCSV(url string) (io.ReadCloser, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get url failed: %w", err)
	}

	return resp.Body, nil
}
