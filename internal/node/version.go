package node

import (
	"errors"
	"strconv"
	"strings"

	"github.com/romycode/amvm/internal"
)

type Version struct {
	Version  string      `json:"version,omitempty"`
	Date     string      `json:"date,omitempty"`
	Files    []string    `json:"files"`
	Npm      string      `json:"npm,omitempty"`
	V8       string      `json:"v8,omitempty"`
	Uv       string      `json:"uv,omitempty"`
	Zlib     string      `json:"zlib,omitempty"`
	Openssl  string      `json:"openssl,omitempty"`
	Modules  string      `json:"modules,omitempty"`
	Lts      interface{} `json:"lts,omitempty"`
	Security bool        `json:"security,omitempty"`
}

func (n Version) IsLts() bool {
	return false != n.Lts
}
func (n Version) MajorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[0])
	return val
}
func (n Version) MinorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[1])
	return val
}
func (n Version) PatchNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[2])
	return val
}
func (n Version) SemverStr() string {
	return n.Version
}
func (n Version) Original() string {
	return n.Version
}
func (n Version) cleanVersion() string {
	return strings.Replace(n.Version, "v", "", 1)
}

type Versions []Version

func (n Versions) Latest() internal.Version {
	version := Version{Version: "v0.0.0"}

	for _, v := range n {
		if v.Version == "" {
			continue
		}

		if version.MajorNum() < v.MajorNum() {
			version = v
		}

		if version.MajorNum() == v.MajorNum() && version.MinorNum() < v.MinorNum() {
			version = v
		}

		if version.MajorNum() == v.MajorNum() && version.MinorNum() == v.MinorNum() && version.PatchNum() < v.PatchNum() {
			version = v
		}
	}

	return version
}
func (n Versions) Lts() internal.Version {
	version := Version{}

	for _, v := range n {
		if v.IsLts() {
			if version.MajorNum() < v.MajorNum() {
				version = v
			}
			if version.MajorNum() == v.MajorNum() && version.MinorNum() < v.MinorNum() {
				version = v
			}
			if version.MajorNum() == v.MajorNum() && version.MinorNum() == v.MinorNum() && version.PatchNum() < v.PatchNum() {
				version = v
			}
		}
	}

	return version
}
func (n Versions) GetVersion(version string) (internal.Version, error) {
	if "latest" == version {
		return n.Latest(), nil
	}

	if "lts" == version {
		return n.Lts(), nil
	}

	if !strings.Contains(version, "v") {
		return Version{}, errors.New("invalid version provided, must start with 'v'")
	}

	ver := strings.Split(version, ".")
	if len(ver) < 3 {
		return Version{}, errors.New("invalid version provided")
	}

	for _, v := range n {
		if v.Version == version {
			return v, nil
		}
	}

	return Version{}, nil
}
