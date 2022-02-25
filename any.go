package scan

import (
	"reflect"

	"github.com/pkg/errors"
)

func copyAny(v any, to any) error {
	if reflect.TypeOf(to).Kind() != reflect.Pointer {
		return errors.Errorf("cannot copy into non-pointer type %T", to)
	}
	toElem := reflect.ValueOf(to).Elem()
	if !toElem.CanSet() {
		return errors.Errorf("cannot set %s", toElem.Type().String())
	}
	rv := reflect.ValueOf(v)
	if !rv.CanConvert(toElem.Type()) {
		return errors.Errorf("cannot convert %T to %s", v, toElem.Type().String())
	}
	crv := rv.Convert(toElem.Type())
	toElem.Set(crv)
	return nil
}
