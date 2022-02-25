package scan

import (
	"reflect"
	"testing"
)

func errWhenNoneExpected(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatalf("error when non expected: %v", err)
}

func noErrWhenErrExpected(t *testing.T, err error) {
	if err != nil {
		return
	}
	t.Fatalf("no-error when err expected")
}

func assertEqual(t *testing.T, want, have any) {
	if reflect.DeepEqual(want, have) {
		return
	}
	t.Fatalf("want %v, have %v", want, have)
}

func TestCopyAny(t *testing.T) {
	var err error
	var n int
	var f float64
	var s string
	type pair struct {
		n int
		s string
	}
	var p pair

	err = copyAny(int(32), &n)
	errWhenNoneExpected(t, err)
	assertEqual(t, 32, n)

	err = copyAny(int(32), &f)
	errWhenNoneExpected(t, err)
	assertEqual(t, 32.0, f)

	err = copyAny(float64(32.34), &f)
	errWhenNoneExpected(t, err)
	assertEqual(t, 32.34, f)

	err = copyAny(float64(32.34), &n)
	errWhenNoneExpected(t, err)
	assertEqual(t, 32, n)

	err = copyAny("hans sausage", &s)
	errWhenNoneExpected(t, err)
	assertEqual(t, "hans sausage", s)

	err = copyAny("hans sausage", &f)
	noErrWhenErrExpected(t, err)

	err = copyAny(pair{2, "xmas"}, &p)
	errWhenNoneExpected(t, err)
	assertEqual(t, pair{2, "xmas"}, p)

	err = copyAny("hans sausage", s)
	noErrWhenErrExpected(t, err)
}
