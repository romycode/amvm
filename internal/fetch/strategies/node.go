package strategies

import (
	"encoding/json"
	"strings"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/http"
)

const (
	NodeJsBaseURL     = "https://nodejs.org"
	nodeJsVersionsURL = "/dist/index.json"
)

type NodeJsFetcherStrategy struct {
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewNodeJsFetcherStrategy(hc *http.DefaultClient, arch, os string) *NodeJsFetcherStrategy {
	return &NodeJsFetcherStrategy{hc, arch, os}
}

func (n NodeJsFetcherStrategy) filterByOsAndArch(versions version.NodeVersions) version.NodeVersions {
	arch := ""
	if "darwin" == n.os {
		arch = "osx-x64-tar"
		if "arm64" == n.arch {
			arch = "osx-arm64-tar"
		}
	}

	if "linux" == n.os {
		arch = "linux-x64"
		if "arm64" == n.arch {
			arch = "linux-arm64"
		}
	}

	filteredVersions := version.NodeVersions{}
	for _, version := range versions {
		if strings.Contains(strings.Join(version.Files, " - "), arch) {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
}
func (n NodeJsFetcherStrategy) Accepts(tool internal.Tool) bool {
	return internal.Node == tool
}

func (n NodeJsFetcherStrategy) Execute() (version.Versions, error) {
	res, err := n.hc.Request("GET", NodeJsBaseURL+nodeJsVersionsURL, "")
	if err != nil {
		return nil, err
	}

	versions := version.NodeVersions{}
	err = json.NewDecoder(res.Body).Decode(&versions)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return n.filterByOsAndArch(versions), nil
}
