package scan

import (
	"fmt"
	"reflect"
	"testing"
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
			in:   "some text {{some evaler}}",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					"some text",
					Evaler{raw: "some evaler"},
				},
			},
			expectPrefix: "some text",
		},
		{
			name: "test",
			in:   "{{some evaler in front of}}some text",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					Evaler{raw: "some evaler in front of"},
					"some text",
				},
			},
			expectPrefix: "",
		},
		{
			name: "test",
			in:   "some text {{some evaler}} and a text behind   ",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					"some text",
					Evaler{raw: "some evaler"},
					"and a text behind",
				},
			},
			expectPrefix: "some text",
		},
		{
			name: "test",
			in:   "   some text direct before{{some evaler}}and direct behind   ",
			fail: false,
			expect: &Template{
				name: "test",
				items: []Item{
					"some text direct before",
					Evaler{raw: "some evaler"},
					"and direct behind",
				},
			},
			expectPrefix: "some text direct before",
		},
		{
			name:         "test",
			in:           "some text {{some evaler}} and {{not correctly closing this}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
		{
			name:         "test",
			in:           "some text {{some evaler}} and {{}}",
			fail:         true,
			expect:       nil,
			expectPrefix: "",
		},
		{
			name:         "test",
			in:           "some text {{some evaler}}{{foo}}",
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
