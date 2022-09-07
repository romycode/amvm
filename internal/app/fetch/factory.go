package fetch

import (
	"fmt"

	"github.com/romycode/amvm/internal/java"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/internal/pnpm"
	"github.com/romycode/amvm/pkg/ui"
)

type Factory struct {
	nf, df, pf, jf internal.Fetcher
}

func NewFactory(nf, df, pf, jf internal.Fetcher) *Factory {
	return &Factory{nf, df, pf, jf}
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

	if java.Java().Value() == tool {
		return ff.jf, nil
	}

	return nil, fmt.Errorf(ui.Colorize("invalid tool", ui.Red))
}
