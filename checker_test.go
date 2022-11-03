package golden_test

import (
	"strings"
	"testing"

	"github.com/tenntenn/golden"
)

func TestChecker_Check(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		want    string
		got     any
		hasDiff bool
	}{
		"string-nodiff":    {"hello", "hello", false},
		"bytes-nodiff":     {"hello", []byte("hello"), false},
		"reader-nodiff":    {"hello", strings.NewReader("hello"), false},
		"json-nodiff":      {"{\"S\":\"hello\"}\n", struct{ S string }{S: "hello"}, false},
		"marshaler-nodiff": {"hello", marshaler("hello"), false},

		"string-diff":    {"Hello", "hello", true},
		"bytes-diff":     {"Hello", []byte("hello"), true},
		"reader-diff":    {"Hello", strings.NewReader("hello"), true},
		"json-diff":      {"{\"S\":\"Hello\"}\n", struct{ S string }{S: "hello"}, true},
		"marshaler-diff": {"Hello", marshaler("hello"), true},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testdata := t.TempDir()

			// only update
			diff1 := golden.New(t, true, testdata, name).Check("_check", tt.want)
			if diff1 != "" {
				t.Error("there are some unexpected differences:", diff1)
			}

			diff2 := golden.New(t, false, testdata, name).Check("_check", tt.got)
			switch {
			case diff2 == "" && tt.hasDiff:
				t.Error("there are any expected differences")
			case diff2 != "" && !tt.hasDiff:
				t.Error("there are some unexpected differences:", diff2)
			}
		})
	}
}
