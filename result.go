package scan

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type ResultItem struct {
	Name  string
	Value any
}

type Result struct {
	Items []ResultItem
}

func (r *Result) Scan(args ...any) error {
	if len(args) > len(r.Items) {
		return errors.Errorf("to many args. want %d at most, got %d", len(r.Items), len(args))
	}
	for i, arg := range args {
		err := copyAny(r.Items[i].Value, arg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Result) Decode(v any) error {
	msa := map[string]any{}
	for _, item := range r.Items {
		msa[item.Name] = item.Value
	}
	bs, err := json.Marshal(msa)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, v)
	return err
}
