package golden_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tenntenn/golden"
)

func TestChecker_Check(t *testing.T) {
	t.Parallel()
	type T struct {
		N int
		M int
	}
	cases := map[string]struct {
		want    string
		got     any
		hasDiff bool
		opts    []cmp.Option
	}{
		"string-nodiff":       {"hello", "hello", false, nil},
		"bytes-nodiff":        {"hello", []byte("hello"), false, nil},
		"reader-nodiff":       {"hello", strings.NewReader("hello"), false, nil},
		"json-nodiff":         {"{\"S\":\"hello\"}\n", struct{ S string }{S: "hello"}, false, nil},
		"marshaler-nodiff":    {"hello", marshaler("hello"), false, nil},
		"ignore-field-nodiff": {`{"N":100, "M":200}`, &T{N: 100, M: 300}, false, []cmp.Option{cmpopts.IgnoreFields(T{}, "M")}},
		"ignore-inner-struct-field-nodiff": {`[{"N":100, "M":200}]`, []*T{{N: 100, M: 300}}, false, []cmp.Option{cmpopts.IgnoreFields(T{}, "M")}},

		"string-diff":    {"Hello", "hello", true, nil},
		"bytes-diff":     {"Hello", []byte("hello"), true, nil},
		"reader-diff":    {"Hello", strings.NewReader("hello"), true, nil},
		"json-diff":      {"{\"S\":\"Hello\"}\n", struct{ S string }{S: "hello"}, true, nil},
		"marshaler-diff": {"Hello", marshaler("hello"), true, nil},
		"ignore-field-diff": {`{"N":100, "M":200}`, &T{N: 100, M: 300}, true, nil},
		"ignore-inner-struct-field-diff": {`[{"N":100, "M":200}]`, []*T{{N: 100, M: 300}}, true, nil},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testdata := t.TempDir()

			// only update
			diff1 := golden.New(t, true, testdata, name, tt.opts...).Check("_check", tt.want)
			if diff1 != "" {
				t.Error("there are some unexpected differences:", diff1)
			}

			diff2 := golden.New(t, false, testdata, name, tt.opts...).Check("_check", tt.got)
			switch {
			case diff2 == "" && tt.hasDiff:
				t.Error("there are any expected differences")
			case diff2 != "" && !tt.hasDiff:
				t.Error("there are some unexpected differences:", diff2)
			}
		})
	}
}
