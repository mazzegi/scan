package scan

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mazzegi/slices"
	"github.com/pkg/errors"
)

func TestParseTemplate(t *testing.T) {
	tests := []struct {
		name         string
		in           string
		fail         bool
		expect       *Template
		expectPrefix string
	}{
		{
			name: "test",
			in:   "",
			fail: false,
			expect: &Template{
				name:  "test",
				items: []Item{},
			},
			expectPrefix: "",
		},
		{
			name: "test",
			in:   "some text {{ some: evaler}}",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					"some text",
					Evaler{
						raw:      "some: evaler",
						name:     "some",
						funcName: "evaler",
					},
				},
			},
			expectPrefix: "some text",
		},
		{
			name: "test",
			in:   "{{some evaler: in front of}}some text",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					Evaler{
						raw:      "some evaler: in front of",
						name:     "some evaler",
						funcName: "in front of",
					},
					"some text",
				},
			},
			expectPrefix: "",
		},
		{
			name: "test",
			in:   "some text {{some :evaler}} and a text behind   ",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					"some text",
					Evaler{
						raw:      "some :evaler",
						name:     "some",
						funcName: "evaler",
					},
					"and a text behind",
				},
			},
			expectPrefix: "some text",
		},
		{
			name: "test",
			in:   "   some text direct before{{some: evaler}}and direct behind   ",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					"some text direct before",
					Evaler{
						raw:      "some: evaler",
						name:     "some",
						funcName: "evaler",
					},
					"and direct behind",
				},
			},
			expectPrefix: "some text direct before",
		},
		{
			name:         "test",
			in:           "some text {{some: evaler}} and {{not correctly: closing this}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
		{
			name:         "test",
			in:           "some text {{some :evaler}} and {{}}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
		{
			name:         "test",
			in:           "some text {{some :evaler}}{{foo: evaler}}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
		{
			name:         "test",
			in:           "some text {{some evaler}}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
		{
			name:         "test",
			in:           "some text {{some: }}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
		{
			name:         "test",
			in:           "some text {{: evaler }}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			res, err := ParseTemplate("test", test.in)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !reflect.DeepEqual(test.expect, res) {
					t.Fatalf("want %v, have %v", test.expect, res)
				}
				if test.expectPrefix != res.Prefix() {
					t.Fatalf("prefix: want %q, have %q", test.expectPrefix, res.Prefix())
				}
				if test.name != res.Name() {
					t.Fatalf("name: want %q, have %q", test.name, res.Name())
				}
			}
		})
	}
}

func TestEvalTemplate(t *testing.T) {
	funcs := BuiltinFuncs()

	tests := []struct {
		template string
		in       string
		fail     bool
		params   map[string]any
	}{
		{
			template: "this one {{name: string}} is here",
			in:       "this one bro",
			fail:     true,
		},
		{
			template: "this one {{name: string}} is here {{tail: int}}",
			in:       "this one bro is here",
			fail:     true,
		},
		{
			template: "the test {{name: string}} will not fail",
			in:       "the test test123 will not fail",
			fail:     false,
			params: map[string]any{
				"name": "test123",
			},
		},
		{
			template: "neither will {{name: string}}, which scans for {{number: int}}",
			in:       "neither will test42$, which scans for 47823",
			fail:     false,
			params: map[string]any{
				"name":   "test42$",
				"number": 47823,
			},
		},
		{
			template: "neither will {{name: string}}, which scans for {{number: float}} and some more",
			in:       "neither will test42$, which scans for 0.7436 and some more",
			fail:     false,
			params: map[string]any{
				"name":   "test42$",
				"number": 0.7436,
			},
		},
		{
			template: "neither will {{name: string}}, which scans for {{flag: bool}} and some more",
			in:       "neither will foo bar baz, which scans for true and some more",
			fail:     false,
			params: map[string]any{
				"name": "foo bar baz",
				"flag": true,
			},
		},
		{
			template: "{{nums: []int}} should be an array of int",
			in:       "42,56, 782 should be an array of int",
			fail:     false,
			params: map[string]any{
				"nums": []int{42, 56, 782},
			},
		},
		{
			template: "{{nums: []int}} should be an array of int followed by an arry of float and the a bool {{floats: []float}} and {{bools: []bool}}",
			in:       "42,56, 782 should be an array of int followed by an arry of float and the a bool 1.3, 7.45, 8976.22, 2e-3 and true, false, true",
			fail:     false,
			params: map[string]any{
				"nums":   []int{42, 56, 782},
				"floats": []float64{1.3, 7.45, 8976.22, 2e-3},
				"bools":  []bool{true, false, true},
			},
		},
		{
			template: "strings should work either: {{heroes: []string}}.",
			in:       "strings should work either: s1, s2, s3, zorro.",
			fail:     false,
			params: map[string]any{
				"heroes": []string{"s1", "s2", "s3", "zorro"},
			},
		},
		{
			template: "this may fail: {{numbers: []int}}.",
			in:       "this may fail: 11,12,boom .",
			fail:     true,
		},
		{
			template: "this may also fail due to a unknown func: {{foo: bar}}.",
			in:       "this may also fail due to a unknown func: bam.",
			fail:     true,
		},
		{
			template: "this one will not match",
			in:       "this is something different",
			fail:     true,
		},
		{
			template: "this one {{name: string}} terminates different",
			in:       "this one bro terminates",
			fail:     true,
		},
		{
			template: "this one {{name: string}} terminates different",
			in:       "this one bro",
			fail:     true,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			tpl, err := ParseTemplate("test", test.template)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			res, err := tpl.Eval(test.in, funcs)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !reflect.DeepEqual(test.params, res.Items) {
					t.Fatalf("want %v, have %v", test.params, res.Items)
				}
			}
		})
	}
}

func TestFuncs(t *testing.T) {
	funcs := BuiltinFuncs()
	funcs.Add("spliturn", func(s string) (any, error) {
		return slices.Convert(strings.Split(s, ":"), slices.TrimSpace)
	})

	tests := []struct {
		template string
		in       string
		fail     bool
		params   map[string]any
	}{
		{
			template: "this urn {{urn: spliturn}} is here",
			in:       "this urn is:a:urn: with is here",
			fail:     false,
			params: map[string]any{
				"urn": []string{"is", "a", "urn", "with"},
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			tpl, err := ParseTemplate("test", test.template)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			res, err := tpl.Eval(test.in, funcs)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !reflect.DeepEqual(test.params, res.Items) {
					t.Fatalf("want %v, have %v", test.params, res.Items)
				}
			}
		})
	}
}

func ptrVal[T any](v T) *T {
	return &v
}

func ptr[T any]() *T {
	var v T
	return &v
}

func TestCapture(t *testing.T) {
	type pair struct {
		first  int
		second int
	}

	funcs := BuiltinFuncs()
	funcs.Add("pair", func(s string) (any, error) {
		ns, err := slices.Convert(strings.Split(s, ":"), slices.ParseInt)
		if err != nil {
			return nil, err
		}
		if len(ns) != 2 {
			return nil, errors.Errorf("invalid arg count for pair")
		}
		return pair{ns[0], ns[1]}, nil
	})

	tests := []struct {
		template    string
		in          string
		captures    []any
		fail        bool
		params      map[string]any
		expCaptures []any
	}{
		{
			template: "a number {{num: int}} comes in",
			in:       "a number 42 comes in",
			captures: []any{ptr[int]()},
			fail:     false,
			params: map[string]any{
				"num": 42,
			},
			expCaptures: []any{ptrVal(42)},
		},
		{
			template: "a string {{str: string}} comes in",
			in:       "a string fortytwo comes in",
			captures: []any{ptr[string]()},
			fail:     false,
			params: map[string]any{
				"str": "fortytwo",
			},
			expCaptures: []any{ptrVal("fortytwo")},
		},
		{
			template: "a string {{str: string}} comes in",
			in:       "a string fortytwo comes in",
			captures: []any{ptr[int]()},
			fail:     true,
		},
		{
			template: "a pair {{foo: pair}} comes in",
			in:       "a pair 1:2 comes in",
			captures: []any{ptr[pair]()},
			fail:     false,
			params: map[string]any{
				"foo": pair{1, 2},
			},
			expCaptures: []any{ptrVal(pair{1, 2})},
		},
		{
			template: "a pair {{foo: pair}} and a float {{num: float}} comes in",
			in:       "a pair 1:2 and a float 46732.123 comes in",
			captures: []any{ptr[pair](), ptr[float64]()},
			fail:     false,
			params: map[string]any{
				"foo": pair{1, 2},
				"num": 46732.123,
			},
			expCaptures: []any{ptrVal(pair{1, 2}), ptrVal(46732.123)},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%02d", i), func(t *testing.T) {
			tpl, err := ParseTemplate("test", test.template)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			res, err := tpl.Eval(test.in, funcs, test.captures...)
			if err != nil {
				if !test.fail {
					t.Fatalf("expect NOT to fail, but got %v", err)
				}
			} else {
				if !reflect.DeepEqual(test.params, res.Items) {
					t.Fatalf("params: want %v, have %v", test.params, res.Items)
				}
				if !reflect.DeepEqual(test.expCaptures, test.captures) {
					t.Fatalf("captures: want %v, have %v", test.expCaptures, test.captures)
				}
			}
		})
	}
}
