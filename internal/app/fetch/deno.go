package fetch

import (
	"encoding/json"
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/pkg/http"
)

const (
	DenoGithubURLTemplate = "https://api.github.com/repos/denoland/deno"
	denoVersionsURL       = "/tags"
)

type DenoFetcher struct {
	hc *http.DefaultClient
}

func NewDenoFetcher(hc *http.DefaultClient) *DenoFetcher {
	return &DenoFetcher{hc: hc}
}

func (n DenoFetcher) Run(flavour string) (internal.Versions, error) {
	res, err := n.hc.Request("GET", fmt.Sprintf(n.hc.URL()+"%s", denoVersionsURL), "")
	if err != nil {
		return nil, err
	}

	versions := deno.Versions{}
	err = json.NewDecoder(res.Body).Decode(&versions)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return versions, nil
}
