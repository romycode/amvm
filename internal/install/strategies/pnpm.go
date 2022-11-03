package strategies

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/http"
	"github.com/romycode/amvm/pkg/ui"
)

type PnpmJsInstallerStrategy struct {
	hc   *http.DefaultClient
	c    internal.Config
	arch string
	os   string
}

func NewPnpmJsInstallerStrategy(hc *http.DefaultClient, c internal.Config, arch, os string) *PnpmJsInstallerStrategy {
	return &PnpmJsInstallerStrategy{hc, c, arch, os}
}

func (n PnpmJsInstallerStrategy) Accepts(tool internal.Tool) bool {
	return internal.Pnpm == tool
}

func (n PnpmJsInstallerStrategy) Execute(ver version.Version) internal.Output {
	target := "linux-x64"
	if "darwin" == n.os {
		target = "macos-x64"
		if "arm64" == n.arch {
			target = "macos-arm64"
		}
	}

	// Pnpm -> https://github.com/pnpm/pnpm/releases/download/v6.32.9/pnpm-linux-arm64
	downloadURL := fmt.Sprintf("https://github.com/pnpm/pnpm/releases/download/%s/pnpm-%s", ver.Original(), target)

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

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	_ = os.MkdirAll(filepath.Join(n.c.VersionsDir, ver.SemverStr(), "bin"), 0755)

	err = file.Write(filepath.Join(n.c.VersionsDir, ver.SemverStr(), "bin", "pnpm"), data)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	return internal.NewOutput(fmt.Sprintf("🔚 Download version: %s 🔚", ver.Original()), ui.Green, 0)
}
