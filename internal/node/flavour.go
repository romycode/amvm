package node

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
	return NodeJs().value == r.value || IoJs().value == r.value
}

func NodeJs() Flavour {
	return Flavour{value: "nodejs"}
}

func IoJs() Flavour {
	return Flavour{value: "iojs"}
}
