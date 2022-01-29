package http

import (
	"io"
	"net/http"
	urlstd "net/url"
	"strings"
)

type Client interface {
	Request(method string, url string, data string) (*http.Response, error)
	URL() string
}

type DefaultClient struct {
	hc  *http.Client
	url string
}

func NewClient(hc *http.Client, url string) *DefaultClient {
	return &DefaultClient{hc: hc, url: url}
}

func (c DefaultClient) URL() string {
	return c.url
}

func (c DefaultClient) Request(method string, url string, data string) (*http.Response, error) {
	fullUrl, err := urlstd.Parse(url)
	if err != nil {
		return &http.Response{}, err
	}

	return c.hc.Do(&http.Request{
		Method: method,
		URL:    fullUrl,
		Body:   io.NopCloser(strings.NewReader(data)),
	})
}
