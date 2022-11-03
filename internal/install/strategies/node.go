package strategies

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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

type NodeJsInstallerStrategy struct {
	hc   *http.DefaultClient
	c    internal.Config
	arch string
	os   string
	tool internal.Tool
}

func NewNodeJsInstallerStrategy(hc *http.DefaultClient, c internal.Config, arch, os string) *NodeJsInstallerStrategy {
	if "darwin" == os {
		arch = "osx-x64-tar"
		if "arm64" == arch {
			arch = "osx-arm64-tar"
		}
	}

	if "linux" == os {
		arch = "x64"
		if "arm64" == arch {
			arch = "linux-arm64"
		}
	}
	return &NodeJsInstallerStrategy{hc, c, arch, os, internal.Node}
}

func (n NodeJsInstallerStrategy) Accepts(tool internal.Tool) bool {
	return tool == n.tool
}

func (n NodeJsInstallerStrategy) Execute(ver version.Version) internal.Output {
	// NodeJs -> https://nodejs.org/dist/v17.3.0/node-v17.3.0-linux-x64.tar.gz
	downloadURL := fmt.Sprintf("https://nodejs.org/dist/%[1]s/node-%[1]s-%[2]s-%[3]s.tar.gz", ver.Original(), n.os, n.arch)
	return n.download(downloadURL, ver, n.c.CacheDir, filepath.Join(n.c.VersionsDir, ver.SemverStr()))
}

func (n NodeJsInstallerStrategy) download(url string, version version.Version, cacheDir string, destDir string) internal.Output {
	spinner := ui.NewSpinner("Downloading version " + version.SemverStr() + "... ")
	spinner.Start()
	defer spinner.Stop()

	res, err := n.hc.Request("GET", url, "")
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	err = file.Write(filepath.Join(cacheDir, version.SemverStr()+".tar.gz"), content)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	gzFile, err := gzip.NewReader(io.NopCloser(bytes.NewReader(content)))
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	tr := tar.NewReader(gzFile)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	dirToMv := cacheDir
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return internal.NewOutput(err.Error(), ui.Red, 1)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if n.c.CacheDir == dirToMv {
				dirToMv = filepath.Join(dirToMv, hdr.Name)
			}
			err := os.MkdirAll(filepath.Join(n.c.CacheDir, hdr.Name), 0755)
			if err != nil {
				return internal.NewOutput(err.Error(), ui.Red, 1)
			}
		case tar.TypeSymlink:
			err := os.Symlink(hdr.Linkname, filepath.Join(n.c.CacheDir, hdr.Name))
			if err != nil {
				return internal.NewOutput(err.Error(), ui.Red, 1)
			}
		default:
			content, err := io.ReadAll(tr)
			if err != nil {
				return internal.NewOutput(err.Error(), ui.Red, 1)
			}

			err = file.Write(filepath.Join(n.c.CacheDir, hdr.Name), content)
			if err != nil {
				return internal.NewOutput(err.Error(), ui.Red, 1)
			}
		}
	}

	err = os.RemoveAll(destDir)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	err = os.Rename(dirToMv, destDir)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}

	return internal.Output{}
}
