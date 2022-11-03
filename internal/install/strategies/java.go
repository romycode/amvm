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

type JavaInstallerStrategy struct {
	hc   *http.DefaultClient
	c    internal.Config
	arch string
	os   string
}

func NewJavaInstallerStrategy(hc *http.DefaultClient, c internal.Config, arch, os string) *JavaInstallerStrategy {
	if "darwin" == os {
		os = "mac"
	}

	if "amd64" == arch {
		arch = "x64"
	}

	if "arm64" == arch {
		arch = "aarch64"
	}

	return &JavaInstallerStrategy{hc, c, arch, os}
}

func (n JavaInstallerStrategy) Accepts(tool internal.Tool) bool {
	return internal.Java == tool
}

func (n JavaInstallerStrategy) Execute(ver version.Version) internal.Output {
	// Java Binary URL -> https://api.adoptium.net/v3/binary/version/jdk-18%2B36/mac/aarch64/jdk/hotspot/normal/eclipse?project=jdk
	downloadURL := fmt.Sprintf("https://api.adoptium.net/v3/binary/version/%s/%s/%s/jdk/hotspot/normal/eclipse?project=jdk", "jdk-"+ver.Original(), n.os, n.arch)
	return n.download(downloadURL, ver, n.c.CacheDir, filepath.Join(n.c.VersionsDir, ver.SemverStr()))
}

func (n JavaInstallerStrategy) download(url string, version version.Version, cacheDir string, destDir string) internal.Output {
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
