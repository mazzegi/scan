package scan

import "github.com/pkg/errors"

type Evaler struct {
	raw string
}

func ParseEvaler(s string) (Evaler, error) {
	if s == "" {
		return Evaler{}, errors.Errorf("empty string")
	}

	e := Evaler{
		raw: s,
	}
	return e, nil
}
