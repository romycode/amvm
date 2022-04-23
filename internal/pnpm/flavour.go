package pnpm

import "errors"

type Flavour struct {
	value string
}

func NewFlavour(val string) (Flavour, error) {
	f := Flavour{value: val}

	if f.IsValid() {
		return f, nil
	}

	return Flavour{}, errors.New("invalid flavour")
}

func (r Flavour) Value() string {
	return r.value
}

func (r Flavour) IsValid() bool {
	return PnpmJs().Value() == r.value || IoJs().Value() == r.value
}

func PnpmJs() Flavour {
	return Flavour{value: "pnpm"}
}

func IoJs() Flavour {
	return Flavour{value: "iojs"}
}
