package fetch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/pkg/http"
)

const (
	NodeJsURLTemplate = "https://%s.org"
	nodeJsVersionsURL = "/dist/index.json"
)

type NodeJsFetcher struct {
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewNodeJsFetcher(hc *http.DefaultClient, arch, os string) *NodeJsFetcher {
	return &NodeJsFetcher{hc, arch, os}
}

func (n NodeJsFetcher) filterByOsAndArch(versions node.Versions) node.Versions {
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

	filteredVersions := node.Versions{}
	for _, version := range versions {
		if strings.Contains(strings.Join(version.Files, " - "), arch) {
			filteredVersions = append(filteredVersions, version)
		}
	}

	return filteredVersions
}

func (n NodeJsFetcher) Run(flavour string) (internal.Versions, error) {
	f, err := node.NewFlavour(flavour)
	if err != nil {
		return nil, err
	}

	url := ""
	if config.NodeFlavour == f {
		url = fmt.Sprintf(n.hc.URL()+"%s", flavour, nodeJsVersionsURL)
	}

	res, err := n.hc.Request("GET", url, "")
	if err != nil {
		return nil, err
	}

	versions := node.Versions{}
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
