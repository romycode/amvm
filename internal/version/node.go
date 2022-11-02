package version

import (
	"errors"
	"strconv"
	"strings"
)

type NodeVersion struct {
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

func (n NodeVersion) IsLts() bool {
	return false != n.Lts
}
func (n NodeVersion) MajorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[0])
	return val
}
func (n NodeVersion) MinorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[1])
	return val
}
func (n NodeVersion) PatchNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[2])
	return val
}
func (n NodeVersion) SemverStr() string {
	return n.Version
}
func (n NodeVersion) Original() string {
	return n.Version
}
func (n NodeVersion) cleanVersion() string {
	return strings.Replace(n.Version, "v", "", 1)
}

type NodeVersions []NodeVersion

func (n NodeVersions) Latest() Version {
	version := NodeVersion{Version: "v0.0.0"}

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
func (n NodeVersions) Lts() Version {
	version := NodeVersion{}

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
func (n NodeVersions) GetVersion(version string) (Version, error) {
	if "latest" == version {
		return n.Latest(), nil
	}

	if "lts" == version {
		return n.Lts(), nil
	}

	if !strings.Contains(version, "v") {
		return NodeVersion{}, errors.New("invalid version provided, must start with 'v'")
	}

	ver := strings.Split(version, ".")
	if len(ver) < 3 {
		return NodeVersion{}, errors.New("invalid version provided")
	}

	for _, v := range n {
		if v.Version == version {
			return v, nil
		}
	}

	return NodeVersion{}, nil
}
