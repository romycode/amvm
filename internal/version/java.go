package version

import (
	"errors"
	"strings"
)

type JavaVersion struct {
	ReleaseName    string
	Major          int    `json:"major"`
	Minor          int    `json:"minor"`
	Patch          int    `json:"patch"`
	Semver         string `json:"semver"`
	Build          int    `json:"build"`
	OpenjdkVersion string `json:"openjdk_version"`
	Security       int    `json:"security"`
}

func (n JavaVersion) IsLts() bool {
	return true
}
func (n JavaVersion) MajorNum() int {
	return n.Major
}
func (n JavaVersion) MinorNum() int {
	return n.Minor
}
func (n JavaVersion) PatchNum() int {
	return n.Patch
}
func (n JavaVersion) SemverStr() string {
	return n.Semver
}
func (n JavaVersion) Original() string {
	return n.OpenjdkVersion
}

type JavaVersions []JavaVersion

func (n JavaVersions) Latest() Version {
	version := JavaVersion{}

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
func (n JavaVersions) Lts() Version {
	version := JavaVersion{}

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
func (n JavaVersions) GetVersion(version string) (Version, error) {
	if "latest" == version {
		return n.Latest(), nil
	}

	if "lts" == version {
		return n.Lts(), nil
	}

	ver := strings.Split(version, ".")
	if len(ver) < 3 {
		return JavaVersion{}, errors.New("invalid version provided")
	}

	for _, v := range n {
		if v.Semver == version {
			return v, nil
		}
	}

	return JavaVersion{}, nil
}
