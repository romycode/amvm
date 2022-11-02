package http

import (
	"net/http"
	urlstd "net/url"
	"strings"
)

type Client interface {
	Request(method string, url string, data string) (*http.Response, error)
}

type DefaultClient struct {
	hc *http.Client
}

func NewClient(hc *http.Client) *DefaultClient {
	return &DefaultClient{hc}
}

func (c DefaultClient) Request(method string, url string, data string) (*http.Response, error) {
	fullUrl, err := urlstd.Parse(url)
	if err != nil {
		return &http.Response{}, err
	}

	req, err := http.NewRequest(method, fullUrl.String(), strings.NewReader(data))
	if err != nil {
		return &http.Response{}, err
	}

	return c.hc.Do(req)
}
