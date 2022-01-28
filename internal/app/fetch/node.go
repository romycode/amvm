package fetch

import (
	"encoding/json"
	"fmt"
	"github.com/romycode/mvm/internal"
	"github.com/romycode/mvm/internal/node"
	"github.com/romycode/mvm/pkg/http"
)

const (
	NodeJsURLTemplate = "https://%s.org"
	nodeJsVersionsURL = "/dist/index.json"
)

type NodeJsFetcher struct {
	hc *http.Client
}

func NewNodeJsFetcher(hc *http.Client) *NodeJsFetcher {
	return &NodeJsFetcher{hc: hc}
}

func (n NodeJsFetcher) Run(flavour string) (internal.Versions, error) {
	f, err := node.NewFlavour(flavour)
	if err != nil {
		return nil, err
	}

	url := ""
	if node.DefaultFlavour == f {
		url = fmt.Sprintf(n.hc.URL+"%s", flavour, nodeJsVersionsURL)
	}
	if node.IoJsFlavour == f {
		url = fmt.Sprintf(n.hc.URL+"%s", flavour, nodeJsVersionsURL)
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

	return versions, nil
}
