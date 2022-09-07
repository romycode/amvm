package java

import (
	"errors"
	"strings"

	"github.com/romycode/amvm/internal"
)

type Version struct {
	ReleaseName    string
	Major          int    `json:"major"`
	Minor          int    `json:"minor"`
	Patch          int    `json:"patch"`
	Semver         string `json:"semver"`
	Build          int    `json:"build"`
	OpenjdkVersion string `json:"openjdk_version"`
	Security       int    `json:"security"`
}

func (n Version) IsLts() bool {
	return true
}
func (n Version) MajorNum() int {
	return n.Major
}
func (n Version) MinorNum() int {
	return n.Minor
}
func (n Version) PatchNum() int {
	return n.Patch
}
func (n Version) SemverStr() string {
	return n.Semver
}
func (n Version) Original() string {
	return n.OpenjdkVersion
}

type Versions []Version

func (n Versions) Latest() internal.Version {
	version := Version{}

	for _, v := range n {
		if version.Major < v.Major {
			version = v
		}

		if version.Major == v.Major && version.Minor < v.Minor {
			version = v
		}

		if version.Major == v.Major && version.Minor == v.Minor && version.Patch < v.Patch {
			version = v
		}
	}

	return version
}
func (n Versions) Lts() internal.Version {
	version := Version{}

	for _, v := range n {
		if v.IsLts() {
			if version.Major < v.Major {
				version = v
			}
			if version.Major == v.Major && version.Minor < v.Minor {
				version = v
			}
			if version.Major == v.Major && version.Minor == v.Minor && version.Patch < v.Patch {
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

	ver := strings.Split(version, ".")
	if len(ver) < 3 {
		return Version{}, errors.New("invalid version provided")
	}

	for _, v := range n {
		if v.Semver == version {
			return v, nil
		}
	}

	return Version{}, nil
}
