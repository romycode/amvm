package strategies

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/http"
	"github.com/romycode/amvm/pkg/ui"
)

type DenoInstallerStrategy struct {
	hc   *http.DefaultClient
	c    internal.Config
	arch string
	os   string
}

func NewDenoInstallerStrategy(hc *http.DefaultClient, c internal.Config, arch, os string) *DenoInstallerStrategy {
	return &DenoInstallerStrategy{hc, c, arch, os}
}

func (n DenoInstallerStrategy) filterByOsAndArch(versions version.DenoVersions) version.DenoVersions {
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
	for _, ver := range versions {
		if arch == ver.Name {
			filteredVersions = append(filteredVersions, ver)
		}
	}

	return filteredVersions
}

func (n DenoInstallerStrategy) Accepts(tool internal.Tool) bool {
	return internal.Deno == tool
}
func (n DenoInstallerStrategy) Execute(ver version.Version) internal.Output {
	target := "x86_64-unknown-linux-gnu"
	if "darwin" == n.os {
		target = "x86_64-apple-darwin"
		if "arm64" == n.arch {
			target = "aarch64-apple-darwin"
		}
	}

	// DenoJs -> https://github.com/denoland/deno/releases/%s/download/deno-%s.zip
	downloadURL := fmt.Sprintf("https://github.com/denoland/deno/releases/download/%s/deno-%s.zip", ver.Original(), target)

	res, err := n.hc.Request("GET", downloadURL, "")
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(res.Body)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)

	}

	if err := file.Extract(content, filepath.Join(n.c.VersionsDir, ver.Original())); err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	return internal.NewOutput(fmt.Sprintf("🔚 Download version: %s 🔚", ver.Original()), ui.Green, 0)
}
