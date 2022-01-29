package internal

type Fetcher interface {
	Run(flavour string) (Versions, error)
}
