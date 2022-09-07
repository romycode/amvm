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
	"github.com/romycode/amvm/internal/app/fetch"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/http"
	"github.com/romycode/amvm/pkg/ui"
)

// InstallCommand command for download required version and save into AMVM_{TOOL}_versions
type InstallCommand struct {
	conf *config.AmvmConfig
	ff   *fetch.Factory
	hc   http.Client
}

// NewInstallCommand return an instance of InstallCommand
func NewInstallCommand(conf *config.AmvmConfig, ff *fetch.Factory, hc http.Client) *InstallCommand {
	return &InstallCommand{conf: conf, ff: ff, hc: hc}
}

// Run get version and download `tar.gz` for save uncompressed into AMVM_{TOOL}_versions
func (i InstallCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm install nodejs v17.3.0", ui.Green, 1)
	}

	system := runtime.GOOS
	arch := runtime.GOARCH

	tool := os.Args[2]
	input := os.Args[3]

	vf, err := i.ff.Build(tool)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	versions, err := vf.Run(tool)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	version, err := versions.GetVersion(input)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	switch tool {
	case config.NodeFlavour.Value():
		if arch == "amd64" {
			arch = "x64"
		}
		// NodeJs -> https://nodejs.org/dist/v17.3.0/node-v17.3.0-linux-x64.tar.gz
		downloadURL := fmt.Sprintf("https://%[1]s.org/dist/%[3]s/%[2]s-%[3]s-%[4]s-%[5]s.tar.gz", tool, strings.Replace(tool, "nodejs", "node", 1), version.Original(), system, arch)
		output, done := downloadNode(i, downloadURL, version)
		if done {
			return output
		}
	case config.JavaFlavour.Value():
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
		output, done := downloadJava(&i, downloadURL, version)
		if done {
			return output
		}
	case config.DenoFlavour.Value():
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

			if err = os.MkdirAll(i.conf.Deno.VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin", 0755); err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}

			content, err := io.ReadAll(f)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}

			err = file.Write(i.conf.Deno.VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin"+string(os.PathSeparator)+zipFile.Name, content)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}

			err = f.Close()
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1)

			}
		}
	case config.PnpmFlavour.Value():
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

		_ = os.MkdirAll(i.conf.Pnpm.VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin", 0755)

		err = file.Write(i.conf.Pnpm.VersionsDir+version.SemverStr()+string(os.PathSeparator)+"bin"+string(os.PathSeparator)+"pnpm", data)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		break
	}

	return NewOutput(fmt.Sprintf("🔚 Download version: %s 🔚", input), ui.Green, 0)
}

func downloadNode(i InstallCommand, downloadURL string, version internal.Version) (Output, bool) {
	spinner := ui.NewSpinner("Downloading version " + version.SemverStr() + "... ")
	spinner.Start()
	defer spinner.Stop()
	res, err := i.hc.Request("GET", downloadURL, "")
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(res.Body)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	gzFile, err := gzip.NewReader(res.Body)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	tr := tar.NewReader(gzFile)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	dirToMv := i.conf.Node.CacheDir
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1), true
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if i.conf.Node.CacheDir == dirToMv {
				dirToMv += hdr.Name
			}
			err := os.MkdirAll(i.conf.Node.CacheDir+hdr.Name, 0755)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		case tar.TypeSymlink:
			err := os.Symlink(hdr.Linkname, i.conf.Node.CacheDir+hdr.Name)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		default:
			content, err := io.ReadAll(tr)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}

			err = file.Write(i.conf.Node.CacheDir+hdr.Name, content)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		}
	}

	err = os.RemoveAll(i.conf.Node.VersionsDir + version.SemverStr())
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	err = os.Rename(dirToMv, i.conf.Node.VersionsDir+version.SemverStr())
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}
	return Output{}, false
}

func downloadJava(i *InstallCommand, downloadURL string, version internal.Version) (Output, bool) {
	spinner := ui.NewSpinner("Downloading version " + version.SemverStr() + "... ")
	spinner.Start()
	defer spinner.Stop()
	res, err := i.hc.Request("GET", downloadURL, "")
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	err = file.Write(i.conf.Java.CacheDir+version.SemverStr()+".tar.gz", content)
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

	dirToMv := i.conf.Java.CacheDir
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1), true
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if i.conf.Java.CacheDir == dirToMv {
				dirToMv += hdr.Name
			}
			err := os.MkdirAll(i.conf.Java.CacheDir+hdr.Name, 0755)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		case tar.TypeSymlink:
			err := os.Symlink(hdr.Linkname, i.conf.Java.CacheDir+hdr.Name)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		default:
			content, err := io.ReadAll(tr)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}

			err = file.Write(i.conf.Java.CacheDir+hdr.Name, content)
			if err != nil {
				return NewOutput(err.Error(), ui.Red, 1), true
			}
		}
	}

	err = os.RemoveAll(i.conf.Java.VersionsDir + version.SemverStr())
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}

	err = os.Rename(dirToMv, i.conf.Java.VersionsDir+version.SemverStr())
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1), true
	}
	return Output{}, false
}
