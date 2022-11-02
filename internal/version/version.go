package version

type Version interface {
	IsLts() bool
	MajorNum() int
	MinorNum() int
	PatchNum() int
	SemverStr() string
	Original() string
}

type Versions interface {
	Latest() Version
	Lts() Version
	GetVersion(version string) (Version, error)
}
