package internal

type Version interface {
	IsLts() bool
	Major() int
	Minor() int
	Patch() int
	Semver() string
}

type Versions interface {
	Latest() Version
	Lts() Version
	GetVersion(version string) (Version, error)
}
