package fetch

import (
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/internal/pnpm"
	"github.com/romycode/amvm/pkg/color"
)

type Factory struct {
	nf, df, pf internal.Fetcher
}

func NewFactory(nf internal.Fetcher, df internal.Fetcher, pf internal.Fetcher) *Factory {
	return &Factory{nf: nf, df: df, pf: pf}
}

func (ff Factory) Build(tool string) (internal.Fetcher, error) {
	if node.NodeJs().Value() == tool {
		return ff.nf, nil
	}

	if deno.DenoJs().Value() == tool {
		return ff.df, nil
	}

	if pnpm.PnpmJs().Value() == tool {
		return ff.pf, nil
	}

	return nil, fmt.Errorf(color.Colorize("invalid tool", color.Red))
}
