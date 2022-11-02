package strategies

import (
	"encoding/json"
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/http"
)

const (
	JavaBaseURL        = "https://api.adoptium.net/v3"
	javaVersionsURL    = "/info/release_versions?architecture=%s&os=%s&heap_size=normal&image_type=jdk&page=0&page_size=50&project=jdk&release_type=ga"
	javaReleaseNameURL = "/info/release_names?os=%s&architecture=%s&version=%s&image_type=jdk&project=jdk&release_type=ga"
)

type JavaFetcherStrategy struct {
	hc   *http.DefaultClient
	arch string
	os   string
}

func NewJavaFetcherStrategy(hc *http.DefaultClient, arch, os string) *JavaFetcherStrategy {
	if "darwin" == os {
		os = "mac"
	}

	if "amd64" == arch {
		arch = "x64"
	}

	if "arm64" == arch {
		arch = "aarch64"
	}

	return &JavaFetcherStrategy{hc, arch, os}
}

func (n JavaFetcherStrategy) Accepts(tool internal.Tool) bool {
	return internal.Java == tool
}

func (n JavaFetcherStrategy) Execute() (version.Versions, error) {
	url := fmt.Sprintf(JavaBaseURL+javaVersionsURL, n.arch, n.os)
	res, err := n.hc.Request("GET", url, "")
	if err != nil {
		return nil, err
	}

	var rawVersions struct {
		List version.JavaVersions `json:"versions"`
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
