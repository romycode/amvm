package deno

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
	return DenoJs().Value() == r.value
}

func DenoJs() Flavour {
	return Flavour{value: "deno"}
}
