package fetch

import (
	"encoding/json"
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/pnpm"
	"github.com/romycode/amvm/pkg/http"
)

const (
	PnpmJsURLTemplate = "https://api.github.com/repos/pnpm/pnpm"
	pnpmJsVersionsURL = "/releases"
)

type PnpmJsFetcher struct {
	hc *http.DefaultClient
}

func NewPnpmJsFetcher(hc *http.DefaultClient) *PnpmJsFetcher {
	return &PnpmJsFetcher{hc: hc}
}

func (n PnpmJsFetcher) Run(flavour string) (internal.Versions, error) {
	_, err := pnpm.NewFlavour(flavour)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(n.hc.URL()+"%s", pnpmJsVersionsURL)
	res, err := n.hc.Request("GET", url, "")
	if err != nil {
		return nil, err
	}

	versions := pnpm.Versions{}
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
