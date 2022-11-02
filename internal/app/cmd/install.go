package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/fetch"
	"github.com/romycode/amvm/internal/fetch/strategies"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/http"
	"github.com/romycode/amvm/pkg/ui"
)

// InstallCommand command for download required version and save into AMVM_{TOOL}_versions
type InstallCommand struct {
	c  *internal.AmvmConfig
	f  *fetch.Fetcher
	hc http.Client
}

// NewInstallCommand return an instance of InstallCommand
func NewInstallCommand(c *internal.AmvmConfig, f *fetch.Fetcher, hc http.Client) *InstallCommand {
	return &InstallCommand{c, f, hc}
}

// Run get version and download `tar.gz` for save uncompressed into AMVM_{TOOL}_versions
func (i InstallCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm install nodejs v17.3.0", ui.Green, 1)
	}

	system := runtime.GOOS
	arch := runtime.GOARCH

	tool := internal.Tool(os.Args[2])
	input := os.Args[3]

	versions, err := i.f.Run(tool)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	version, err := versions.GetVersion(input)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	switch tool {
	case internal.Node:
		if arch == "amd64" {
			arch = "x64"
		}
		// NodeJs -> https://nodejs.org/dist/v17.3.0/node-v17.3.0-linux-x64.tar.gz
		downloadURL := fmt.Sprintf("%s/dist/%[3]s/%[2]s-%[3]s-%[4]s-%[5]s.tar.gz", strategies.NodeJsBaseURL, string(tool), version.Original(), system, arch)
		output, done := i.download(downloadURL, internal.Node, version, i.c.Tools[internal.Node].CacheDir, filepath.Join(i.c.Tools[internal.Node].VersionsDir, version.SemverStr()))
		if done {
			return output
		}
	case internal.Java:
		osTarget := runtime.GOOS
		if "darwin" == osTarget {
			osTarget = "mac"
		}
		if arch == "amd64" {
			arch = "x64"
		}
		if arch == "arm64" {
			arch = "aarch64"
		}

		// Java Binary URL -> https://api.adoptium.net/v3/binary/version/jdk-18%2B36/mac/aarch64/jdk/hotspot/normal/eclipse?project=jdk
		downloadURL := fmt.Sprintf("https://api.adoptium.net/v3/binary/version/%s/%s/%s/jdk/hotspot/normal/eclipse?project=jdk", "jdk-"+version.Original(), osTarget, arch)
		output, done := i.download(downloadURL, internal.Java, version, i.c.Tools[internal.Java].CacheDir, i.c.Tools[internal.Java].VersionsDir+version.SemverStr())
		if done {
			return output
		}
	case internal.Deno:
		target := "x86_64-unknown-linux-gnu"
		if "darwin" == system {
			target = "x86_64-apple-darwin"
			if "arm64" == arch {
				target = "aarch64-apple-darwin"
			}
		}

		// DenoJs -> https://github.com/denoland/deno/releases/%s/download/deno-%s.zip
		downloadURL := fmt.Sprintf("https://github.com/denoland/deno/releases/download/%s/deno-%s.zip", version.Original(), target)

		res, err := i.hc.Request("GET", downloadURL, "")
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
		}(res.Body)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		for _, zipFile := range zipReader.File {
			f, err := zipFile.Open()
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}

			if err = os.MkdirAll(i.c.Tools[internal.Deno].VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin", 0755); err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}

			content, err := io.ReadAll(f)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}

			err = file.Write(i.c.Tools[internal.Deno].VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin"+string(os.PathSeparator)+zipFile.Name, content)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}

			err = f.Close()
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}
		}
	case internal.Pnpm:
		target := "linux-x64"
		if "darwin" == system {
			target = "macos-x64"
			if "arm64" == arch {
				target = "macos-arm64"
			}
		}

		// Pnpm -> https://github.com/pnpm/pnpm/releases/download/v6.32.9/pnpm-linux-arm64
		downloadURL := fmt.Sprintf("https://github.com/pnpm/pnpm/releases/download/%s/pnpm-%s", version.Original(), target)

		res, err := i.hc.Request("GET", downloadURL, "")
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
		}(res.Body)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		_ = os.MkdirAll(i.c.Tools[internal.Pnpm].VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin", 0755)

		err = file.Write(i.c.Tools[internal.Pnpm].VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin"+string(os.PathSeparator)+"pnpm", data)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		break
	}

	return NewOutput(fmt.Sprintf("🔚 Download version: %s 🔚", input), ui.Green, 0)
}

func (i InstallCommand) download(url string, tool internal.Tool, version version.Version, cacheDir string, destDir string) (Output, bool) {
	spinner := ui.NewSpinner("Downloading version " + version.SemverStr() + "... ")
	spinner.Start()
	defer spinner.Stop()

	res, err := i.hc.Request("GET", url, "")
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	err = file.Write(filepath.Join(cacheDir, version.SemverStr()+".tar.gz"), content)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	gzFile, err := gzip.NewReader(io.NopCloser(bytes.NewReader(content)))
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	tr := tar.NewReader(gzFile)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	dirToMv := cacheDir
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return NewOutput(err.Error(), ui.Red, 1), true
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if i.c.Tools[tool].CacheDir == dirToMv {
				dirToMv = filepath.Join(dirToMv, hdr.Name)
			}
			err := os.MkdirAll(filepath.Join(i.c.Tools[tool].CacheDir, hdr.Name), 0755)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		case tar.TypeSymlink:
			err := os.Symlink(hdr.Linkname, filepath.Join(i.c.Tools[tool].CacheDir, hdr.Name))
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		default:
			content, err := io.ReadAll(tr)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}

			err = file.Write(filepath.Join(i.c.Tools[tool].CacheDir, hdr.Name), content)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		}
	}

	err = os.RemoveAll(destDir)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	err = os.Rename(dirToMv, destDir)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}
	return Output{}, false
}
