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
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewPnpmJsFetcher(hc *http.DefaultClient, arch, os string) *PnpmJsFetcher {
	return &PnpmJsFetcher{hc, arch, os}
}

func (n PnpmJsFetcher) filterByOsAndArch(versions pnpm.Versions) pnpm.Versions {
	arch := ""
	if "darwin" == n.os {
		arch = "pnpm-macos-x64"
		if "arm64" == n.arch {
			arch = "pnpm-macos-arm64"
		}
	}

	if "Linux" == n.os {
		arch = "pnpm-linux-x64"
		if "arm64" == n.arch {
			arch = "pnpm-linux-arm64"
		}
	}

	filteredVersions := pnpm.Versions{}
	for _, version := range versions {
		if arch == version.Name {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
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
