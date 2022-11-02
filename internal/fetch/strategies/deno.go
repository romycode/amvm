package strategies

import (
	"encoding/json"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/http"
)

const (
	DenoGithubBaseURL = "https://api.github.com/repos/denoland/deno"
	denoVersionsURL   = "/releases"
)

type DenoFetcherStrategy struct {
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewDenoFetcherStrategy(hc *http.DefaultClient, arch, os string) *DenoFetcherStrategy {
	return &DenoFetcherStrategy{hc, arch, os}
}

func (n DenoFetcherStrategy) filterByOsAndArch(versions version.DenoVersions) version.DenoVersions {
	arch := ""
	if "darwin" == n.os {
		arch = "deno-x86_64-apple-darwin.zip"
		if "arm64" == n.arch {
			arch = "deno-aarch64-apple-darwin.zip"
		}
	}

	if "Linux" == n.os {
		arch = "deno-x86_64-unknown-linux-gnu.zip"
	}

	filteredVersions := version.DenoVersions{}
	for _, version := range versions {
		if arch == version.Name {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
}

func (n DenoFetcherStrategy) Accepts(tool internal.Tool) bool {
	return internal.Deno == tool
}
func (n DenoFetcherStrategy) Execute() (version.Versions, error) {
	res, err := n.hc.Request("GET", DenoGithubBaseURL+denoVersionsURL, "")
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
