package scan

import (
	"strings"

	"github.com/pkg/errors"
)

type Evaler struct {
	raw      string
	name     string
	funcName string
}

func ParseEvaler(s string) (Evaler, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Evaler{}, errors.Errorf("empty string")
	}
	name, funcName, ok := strings.Cut(s, ":")
	if !ok {
		return Evaler{}, errors.Errorf("invalid syntax. not in form <name:funcName>")
	}
	name = strings.TrimSpace(name)
	funcName = strings.TrimSpace(funcName)
	if name == "" {
		return Evaler{}, errors.Errorf("empty name")
	}
	if funcName == "" {
		return Evaler{}, errors.Errorf("empty func-name")
	}

	e := Evaler{
		raw:      s,
		name:     name,
		funcName: funcName,
	}
	return e, nil
}
