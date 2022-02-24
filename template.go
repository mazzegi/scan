package scan

import (
	"strings"

	"github.com/pkg/errors"
)

type Result struct {
	items map[string]any
}

type Item interface{}

type Template struct {
	name  string
	items []Item
}

func (t *Template) Name() string {
	return t.name
}

func (t *Template) Prefix() string {
	if len(t.items) == 0 {
		return ""
	}
	if s, ok := t.items[0].(string); ok {
		return s
	}
	return ""
}

//
func (t *Template) Eval(s string, funcs Funcs) (*Result, error) {
	res := &Result{
		items: map[string]any{},
	}
	s = strings.TrimSpace(s)
	var pos int = 0

	eatWhite := func() int {
		var eaten int
		for {
			if s[pos] != ' ' {
				return eaten
			}
			pos++
			if pos >= len(s) {
				return eaten
			}
			eaten++
		}
	}

	for i, item := range t.items {
		eatWhite()
		if pos >= len(s) {
			return nil, errors.Errorf("EOF")
		}
		switch item := item.(type) {
		case string:
			if !strings.HasPrefix(s[pos:], item) {
				return nil, errors.Errorf("no match for string %q", item)
			}
			pos += len(item)
		case Evaler:
			var es string
			if i == len(t.items)-1 {
				es = s[pos:]
			} else {
				//peek next string
				next, ok := t.items[i+1].(string)
				if !ok {
					return nil, errors.Errorf("next is not a string")
				}
				nextIdx := strings.Index(s[pos:], next)
				if nextIdx < 0 {
					return nil, errors.Errorf("no match for next %q", next)
				}
				es = s[pos : pos+nextIdx]
			}
			es = strings.TrimSpace(es)

			v, err := item.Eval(es, funcs)
			if err != nil {
				return nil, errors.Wrapf(err, "eval %q", es)
			}
			res.items[item.name] = v
			pos += len(es)
		}
	}

	return res, nil
}
