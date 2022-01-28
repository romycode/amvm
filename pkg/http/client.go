package http

import (
	"io"
	"net/http"
	urlstd "net/url"
	"strings"
)

type Client struct {
	hc  *http.Client
	URL string
}

func NewClient(hc *http.Client, url string) *Client {
	return &Client{hc: hc, URL: url}
}

func (f Client) Request(method string, url string, data string) (*http.Response, error) {
	fullUrl, err := urlstd.Parse(url)
	if err != nil {
		return &http.Response{}, err
	}

	return f.hc.Do(&http.Request{
		Method: method,
		URL:    fullUrl,
		Body:   io.NopCloser(strings.NewReader(data)),
	})
}
