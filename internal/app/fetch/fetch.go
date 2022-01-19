package fetch

import "github.com/romycode/mvm/internal"

type Fetcher interface {
	Run(flavour string) internal.Versions
}
