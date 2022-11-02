package fetch

import (
	"errors"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
)

type Strategy interface {
	Execute() (version.Versions, error)
	Accepts(tool internal.Tool) bool
}

type Fetcher struct {
	strategies []Strategy
}

func NewFetcher(strategies []Strategy) *Fetcher {
	return &Fetcher{strategies: strategies}
}

func (r Fetcher) Run(tool internal.Tool) (version.Versions, error) {
	for _, strategy := range r.strategies {
		if strategy.Accepts(tool) {
			return strategy.Execute()
		}
	}

	return nil, errors.New("failed fetching versions")
}
