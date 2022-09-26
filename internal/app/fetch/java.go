package fetch

import (
	"encoding/json"
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/java"
	"github.com/romycode/amvm/pkg/http"
)

const (
	JavaURLApi         = "https://api.adoptium.net/v3"
	javaVersionsURL    = "/info/release_versions?architecture=%s&os=%s&heap_size=normal&image_type=jdk&page=0&page_size=50&project=jdk&release_type=ga"
	javaReleaseNameURL = "/info/release_names?os=%s&architecture=%s&version=%s&image_type=jdk&project=jdk&release_type=ga"
)

type JavaFetcher struct {
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewJavaFetcher(hc *http.DefaultClient, arch, os string) *JavaFetcher {
	if "linux" == os {
		os = "linux"
	}

	if "darwin" == os {
		os = "mac"
	}

	if "amd64" == arch {
		arch = "x64"
	}

	if "arm64" == arch {
		arch = "aarch64"
	}

	return &JavaFetcher{hc, arch, os}
}

func (n JavaFetcher) Run(_ string) (internal.Versions, error) {
	versionsURL := n.hc.URL() + javaVersionsURL
	res, err := n.hc.Request("GET", fmt.Sprintf(versionsURL, n.arch, n.os), "")
	if err != nil {
		return nil, err
	}

	var rawVersions struct {
		List java.Versions `json:"versions"`
	}
	err = json.NewDecoder(res.Body).Decode(&rawVersions)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return rawVersions.List, nil
}
