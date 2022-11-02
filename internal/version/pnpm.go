package version

import (
	"errors"
	"strconv"
	"strings"
)

type PnpmVersion struct {
	Name   string `json:"tag_name"`
	Assets []struct {
		Name string `json:"name"`
	} `json:"assets"`
}

func (n PnpmVersion) IsLts() bool {
	return false
}
func (n PnpmVersion) MajorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[0])
	return val
}
func (n PnpmVersion) MinorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[1])
	return val
}
func (n PnpmVersion) PatchNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[2])
	return val
}
func (n PnpmVersion) SemverStr() string {
	return n.Name
}
func (n PnpmVersion) Original() string {
	return n.Name
}
func (n PnpmVersion) cleanVersion() string {
	return strings.Replace(n.Name, "v", "", 1)
}

type PnpmVersions []PnpmVersion

func (n PnpmVersions) Latest() Version {
	version := PnpmVersion{Name: "v0.0.0"}

	for _, v := range n {
		if v.Name == "" {
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
func (n PnpmVersions) Lts() Version {
	version := PnpmVersion{}

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
func (n PnpmVersions) GetVersion(version string) (Version, error) {
	if "latest" == version {
		return n.Latest(), nil
	}

	if "lts" == version {
		return n.Lts(), nil
	}

	if !strings.Contains(version, "v") {
		return PnpmVersion{}, errors.New("invalid version provided, must start with 'v'")
	}

	ver := strings.Split(version, ".")
	if len(ver) < 3 {
		return PnpmVersion{}, errors.New("invalid version provided")
	}

	for _, v := range n {
		if v.Name == version {
			return v, nil
		}
	}

	return PnpmVersion{}, nil
}
