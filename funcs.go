package scan

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type EvalFunc func(s string) (any, error)

type Funcs map[string]EvalFunc

func (fs Funcs) Add(name string, fnc EvalFunc) {
	fs[name] = fnc
}

func builtinFuncs() Funcs {
	fs := Funcs{}
	fs["string"] = func(s string) (any, error) {
		return s, nil
	}
	fs["int"] = func(s string) (any, error) {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-int %q", s)
		}
		return n, nil
	}
	fs["float"] = func(s string) (any, error) {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-float %q", s)
		}
		return f, nil
	}
	fs["bool"] = func(s string) (any, error) {
		t, err := strconv.ParseBool(s)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-bool %q", s)
		}
		return t, nil
	}
	fs["[]string"] = func(s string) (any, error) {
		return strings.Split(s, ","), nil
	}

	return fs
}
