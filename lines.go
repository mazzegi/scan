package scan

import (
	"bufio"
	"io"
	"strings"

	"github.com/pkg/errors"
)

func Lines[T any](pattern string, funcs Funcs, r io.Reader) ([]T, error) {
	tpl, err := ParseTemplate("lines", pattern)
	if err != nil {
		return nil, errors.Wrap(err, "parse template")
	}
	var ts []T
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		if ln == "" {
			continue
		}
		res, err := tpl.Eval(ln, funcs)
		if err != nil {
			return nil, errors.Wrapf(err, "eval line %q", ln)
		}
		var t T
		err = res.Decode(&t)
		if err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		ts = append(ts, t)
	}
	return ts, nil
}
