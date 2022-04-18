package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/pkg/color"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/http"
)

// InstallCommand command for download required version and save into AMVM_{TOOL}_versions
type InstallCommand struct {
	conf *config.AmvmConfig
	nf   internal.Fetcher
	df   internal.Fetcher
	nhc  http.Client
	dhc  http.Client
}

// NewInstallCommand return an instance of InstallCommand
func NewInstallCommand(conf *config.AmvmConfig, nf internal.Fetcher, df internal.Fetcher, nhc http.Client, dhc http.Client) *InstallCommand {
	return &InstallCommand{conf: conf, nf: nf, df: df, nhc: nhc, dhc: dhc}
}

// Run get version and download `tar.gz` for save uncompressed into AMVM_{TOOL}_versions
func (i InstallCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm install nodejs v17.3.0", 1)
	}

	system := runtime.GOOS
	arch := runtime.GOARCH

	tool := os.Args[2]
	_, notNodeTool := node.NewFlavour(tool)
	_, notDenoTool := deno.NewFlavour(tool)
	if notNodeTool != nil && notDenoTool != nil {
		message := notNodeTool.Error()
		if notDenoTool != nil {
			message = notDenoTool.Error()
		}
		return NewOutput(message, 1)
	}

	input := os.Args[3]
	if config.IoJsFlavour.Value() == tool || config.DefaultNodeJsFlavour.Value() == tool {
		versions, err := i.nf.Run(tool)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		version, err := versions.GetVersion(input)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		// https://nodejs.org/dist/v17.3.0/node-v17.3.0-linux-arm64.tar.gz
		// https://iojs.org/dist/v3.3.1/iojs-v3.3.1-linux-x64.tar.gz
		downloadURL := fmt.Sprintf(i.nhc.URL()+"/dist/%[3]s/%[2]s-%[3]s-%[4]s-%[5]s.tar.gz", tool, strings.Replace(tool, "nodejs", "node", 1), version.Semver(), system, arch)
		res, err := i.nhc.Request("GET", downloadURL, "")
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
		}(res.Body)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

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
				return NewOutput(err.Error(), 1)
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

	if config.DefaultDenoJsFlavour.Value() == tool {
		versions, err := i.df.Run(tool)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		version, err := versions.GetVersion(input)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		target := "x86_64-unknown-linux-gnu"
		if "darwin" == system {
			target = "x86_64-apple-darwin"
			if "amd64" == arch {
				target = "aarch64-apple-darwin"
			}
		}

		// https://github.com/denoland/deno/releases/%s/download/deno-%s.zip
		downloadURL := fmt.Sprintf("https://github.com/denoland/deno/releases/download/%s/deno-%s.zip", input, target)
		res, err := i.dhc.Request("GET", downloadURL, "")
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
		}(res.Body)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		for _, zipFile := range zipReader.File {
			f, err := zipFile.Open()
			if err != nil {
				return NewOutput(err.Error(), 1)
			}

			if err = os.MkdirAll(i.conf.Deno.VersionsDir+version.Semver(), 0755); err != nil {
				return NewOutput(err.Error(), 1)
			}

			content, err := io.ReadAll(f)
			if err != nil {
				return NewOutput(err.Error(), 1)
			}

			err = file.Write(i.conf.Deno.VersionsDir+version.Semver()+string(os.PathSeparator)+zipFile.Name, content)
			if err != nil {
				return NewOutput(err.Error(), 1)
			}

			err = f.Close()
			if err != nil {
				return NewOutput(err.Error(), 1)
			}
		}
	}

	return NewOutput(color.Colorize(fmt.Sprintf("🔚 Download version: %s 🔚", input), color.Green), 0)
}
