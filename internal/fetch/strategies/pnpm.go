package strategies

import (
	"encoding/json"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/http"
)

const (
	PnpmJsBaseURL     = "https://api.github.com/repos/pnpm/pnpm"
	pnpmJsVersionsURL = "/releases"
)

type PnpmJsFetcherStrategy struct {
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewPnpmJsFetcherStrategy(hc *http.DefaultClient, arch, os string) *PnpmJsFetcherStrategy {
	return &PnpmJsFetcherStrategy{hc, arch, os}
}

func (n PnpmJsFetcherStrategy) filterByOsAndArch(versions version.DenoVersions) version.DenoVersions {
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

	filteredVersions := version.DenoVersions{}
	for _, version := range versions {
		if arch == version.Name {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
}

func (n PnpmJsFetcherStrategy) Accepts(tool internal.Tool) bool {
	return internal.Pnpm == tool
}

func (n PnpmJsFetcherStrategy) Execute() (version.Versions, error) {
	res, err := n.hc.Request("GET", PnpmJsBaseURL+pnpmJsVersionsURL, "")
	if err != nil {
		return nil, err
	}

	versions := version.DenoVersions{}
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
