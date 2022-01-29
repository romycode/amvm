package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/romycode/mvm/internal"
	"github.com/romycode/mvm/internal/config"
	"github.com/romycode/mvm/internal/node"
	"github.com/romycode/mvm/pkg/color"
	"github.com/romycode/mvm/pkg/file"
	"github.com/romycode/mvm/pkg/http"
)

// InstallCommand command for download required version and save into MVM_{TOOL}_versions
type InstallCommand struct {
	conf *config.MvmConfig
	nf   internal.Fetcher
	hc   http.Client
}

// NewInstallCommand return an instance of InstallCommand
func NewInstallCommand(conf *config.MvmConfig, nf internal.Fetcher, hc http.Client) *InstallCommand {
	return &InstallCommand{conf: conf, nf: nf, hc: hc}
}

// Run get version and download tar.gz for save uncompressed into MVM_{TOOL}_versions
func (i InstallCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: mvm install nodejs v17.3.0", 1)
	}

	system := runtime.GOOS
	arch := runtime.GOARCH
	if "amd64" == arch {
		arch = "x64"
	}

	tool, err := node.NewFlavour(os.Args[2])
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	input := os.Args[3]
	if config.IoJsFlavour == tool || config.DefaultFlavour == tool {
		versions, err := i.nf.Run(tool.Value())
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		version, err := versions.GetVersion(input)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		downloadURL := fmt.Sprintf(i.hc.URL()+"/dist/%[2]s/node-%[2]s-%[3]s-%[4]s.tar.gz", tool.Value(), version.Semver(), system, arch)
		res, err := i.hc.Request("GET", downloadURL, "")
		if err != nil {
			return NewOutput(err.Error(), 1)
		}
		defer res.Body.Close()

		gzFile, err := gzip.NewReader(res.Body)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		tr := tar.NewReader(gzFile)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		dirToMv := i.conf.Node.CacheDir
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
			}

			switch hdr.Typeflag {
			case tar.TypeDir:
				if i.conf.Node.CacheDir == dirToMv {
					dirToMv += hdr.Name
				}
				err := os.MkdirAll(i.conf.Node.CacheDir+hdr.Name, 0755)
				if err != nil {
					return NewOutput(err.Error(), 1)
				}
			case tar.TypeSymlink:
				err := os.Symlink(hdr.Linkname, i.conf.Node.CacheDir+hdr.Name)
				if err != nil {
					return NewOutput(err.Error(), 1)
				}
			default:
				content, err := io.ReadAll(tr)
				if err != nil {
					return NewOutput(err.Error(), 1)
				}

				err = file.Write(i.conf.Node.CacheDir+hdr.Name, content)
				if err != nil {
					return NewOutput(err.Error(), 1)
				}
			}
		}

		err = os.RemoveAll(i.conf.Node.VersionsDir + version.Semver())
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		err = os.Rename(dirToMv, i.conf.Node.VersionsDir+version.Semver())
		if err != nil {
			return NewOutput(err.Error(), 1)
		}
	}

	return NewOutput(color.Colorize(fmt.Sprintf("🔚 Download version: %s 🔚", input), color.Green), 0)
}
