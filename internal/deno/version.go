package deno

import (
	"errors"
	"strconv"
	"strings"

	"github.com/romycode/amvm/internal"
)

type Version struct {
	Name string `json:"name"`
}

func (n Version) IsLts() bool {
	return true
}
func (n Version) Major() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[0])
	return val
}
func (n Version) Minor() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[1])
	return val
}
func (n Version) Patch() int {
	val, _ := strconv.Atoi(strings.Split(n.cleanVersion(), ".")[2])
	return val
}
func (n Version) Semver() string {
	return n.cleanVersion()
}
func (n Version) Original() string {
	return n.Name
}
func (n Version) cleanVersion() string {
	return strings.Replace(n.Name, "v", "", 1)
}

type Versions []Version

func (n Versions) Latest() internal.Version {
	version := Version{Name: "v0.0.0"}

	for _, v := range n {
		if version.Major() < v.Major() {
			version = v
		}

		if version.Major() == v.Major() && version.Minor() < v.Minor() {
			version = v
		}

		if version.Major() == v.Major() && version.Minor() == v.Minor() && version.Patch() < v.Patch() {
			version = v
		}
	}

	return version
}
func (n Versions) Lts() internal.Version {
	version := Version{}

	for _, v := range n {
		if v.IsLts() {
			if version.Major() < v.Major() {
				version = v
			}
			if version.Major() == v.Major() && version.Minor() < v.Minor() {
				version = v
			}
			if version.Major() == v.Major() && version.Minor() == v.Minor() && version.Patch() < v.Patch() {
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
		if v.Name == version {
			return v, nil
		}
	}

	return Version{}, nil
}
