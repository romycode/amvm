package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
	"github.com/romycode/mvm/internal/node"
	"github.com/romycode/mvm/pkg/color"
)

type InstallCommand struct {
	conf *config.MvmConfig
	nf   fetch.Fetcher
}

func NewInstallCommand(conf *config.MvmConfig, nf fetch.Fetcher) *InstallCommand {
	return &InstallCommand{conf: conf, nf: nf}
}

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

	if config.IoJs == tool || config.NodeJs == tool {
		versions, err := i.nf.Run(tool.Value())
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		version, err := versions.GetVersion(input)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		downloadURL := fmt.Sprintf("https://nodejs.org/dist/%[1]s/node-%[1]s-%[2]s-%[3]s.tar.gz", version.Semver(), system, arch)
		if config.IoJs == tool {
			downloadURL = fmt.Sprintf("https://iojs.org/dist/%[1]s/iojs-%[1]s-%[2]s-%[3]s.tar.gz", version.Semver(), system, arch)
		}

		res, err := http.Get(downloadURL)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)

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
				f, err := os.OpenFile(i.conf.Node.CacheDir+hdr.Name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					return NewOutput(err.Error(), 1)
				}

				_, err = io.Copy(f, tr)
				if err != nil {
					return NewOutput(err.Error(), 1)
				}

				_ = f.Close()
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
