package scan

import (
	"strconv"
	"strings"

	"github.com/mazzegi/slices"
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
		return int(n), err
	}
	fs["float"] = func(s string) (any, error) {
		return strconv.ParseFloat(s, 64)
	}
	fs["bool"] = func(s string) (any, error) {
		return strconv.ParseBool(s)
	}
	fs["[]string"] = func(s string) (any, error) {
		return slices.Convert(strings.Split(s, ","), slices.TrimSpace)
	}
	fs["[]int"] = func(s string) (any, error) {
		return slices.Convert(strings.Split(s, ","), slices.ParseInt)
	}
	fs["[]float"] = func(s string) (any, error) {
		return slices.Convert(strings.Split(s, ","), slices.ParseFloat)
	}
	fs["[]bool"] = func(s string) (any, error) {
		return slices.Convert(strings.Split(s, ","), slices.ParseBool)
	}

	return fs
}
