package version

import (
	"errors"
	"strconv"
	"strings"
)

type DenoVersion struct {
	Name string `json:"name"`
}

func (n DenoVersion) IsLts() bool {
	return false
}
func (n DenoVersion) MajorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[0])
	return val
}
func (n DenoVersion) MinorNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[1])
	return val
}
func (n DenoVersion) PatchNum() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[2])
	return val
}
func (n DenoVersion) SemverStr() string {
	return n.cleanVersion()
}
func (n DenoVersion) Original() string {
	return n.Name
}
func (n DenoVersion) cleanVersion() string {
	return strings.Replace(n.Name, "v", "", 1)
}

type DenoVersions []DenoVersion

func (n DenoVersions) Latest() Version {
	version := DenoVersion{Name: "v0.0.0"}

	for _, v := range n {
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
func (n DenoVersions) Lts() Version {
	version := DenoVersion{}

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
func (n DenoVersions) GetVersion(version string) (Version, error) {
	if "latest" == version {
		return n.Latest(), nil
	}

	if "lts" == version {
		return n.Lts(), nil
	}

	if !strings.Contains(version, "v") {
		return DenoVersion{}, errors.New("invalid version provided, must start with 'v'")
	}

	ver := strings.Split(version, ".")
	if len(ver) < 3 {
		return DenoVersion{}, errors.New("invalid version provided")
	}

	for _, v := range n {
		if v.Name == version {
			return v, nil
		}
	}

	return DenoVersion{}, nil
}
