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
	denoVersionsURL       = "/releases"
)

type DenoFetcher struct {
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewDenoFetcher(hc *http.DefaultClient, arch, os string) *DenoFetcher {
	return &DenoFetcher{hc, arch, os}
}

func (n DenoFetcher) filterByOsAndArch(versions deno.Versions) deno.Versions {
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

	filteredVersions := deno.Versions{}
	for _, version := range versions {
		if arch == version.Name {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
}

func (n DenoFetcher) Run(_ string) (internal.Versions, error) {
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
