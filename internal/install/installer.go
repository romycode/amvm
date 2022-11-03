package install

import (
	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/ui"
)

type Strategy interface {
	Execute(ver version.Version) internal.Output
	Accepts(tool internal.Tool) bool
}

type Installer struct {
	strategies []Strategy
}

func NewInstaller(strategies []Strategy) *Installer {
	return &Installer{strategies: strategies}
}

func (r Installer) Run(tool internal.Tool, ver version.Version) internal.Output {
	for _, strategy := range r.strategies {
		if strategy.Accepts(tool) {
			return strategy.Execute(ver)
		}
	}

	return internal.NewOutput("no installer found", ui.Red, 1)
}
