package fetch

import (
	"encoding/json"
	"net/http"

	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/node"
)

const (
	NodeJsURL         = "https://nodejs.org"
	NodeJsVersionsURL = "/dist/index.json"

	IoJsURL         = "https://iojs.org"
	IoJsVersionsURL = "/dist/index.json"
)

func NodeVersions(flavour node.Flavour) (node.Versions, error) {
	url := ""

	if config.NodeJs == flavour {
		url = NodeJsURL + NodeJsVersionsURL
	}

	if config.IoJs == flavour {
		url = IoJsURL + IoJsVersionsURL
	}

	res, err := http.Get(url)
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
