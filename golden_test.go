package golden_test

import (
	"flag"
	"strings"
	"testing"

	"github.com/tenntenn/golden"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

type marshaler string

func (m marshaler) MarshalText() (text []byte, err error) {
	return []byte(m), nil
}

func TestDiff(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		want    string
		got     any
		hasDiff bool
	}{
		"string-nodiff":       {"hello", "hello", false},
		"bytes-nodiff":        {"hello", []byte("hello"), false},
		"reader-nodiff":       {"hello", strings.NewReader("hello"), false},
		"json-nodiff":         {"{\"S\":\"hello\"}\n", struct{ S string }{S: "hello"}, false},
		"marshaler-nodiff":    {"hello", marshaler("hello"), false},
		"empty-nodiff":        {"", "", false},
		"number-nodiff":       {"3", "3", false},
		"number-start-nodiff": {"3 bytes", "3 bytes", false},

		"string-diff":       {"Hello", "hello", true},
		"bytes-diff":        {"Hello", []byte("hello"), true},
		"reader-diff":       {"Hello", strings.NewReader("hello"), true},
		"json-diff":         {"{\"S\":\"Hello\"}\n", struct{ S string }{S: "hello"}, true},
		"marshaler-diff":    {"Hello", marshaler("hello"), true},
		"number-diff":       {"3", "4", true},
		"number-start-diff": {"3 bytes", "4 bytes", true},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testdata := t.TempDir()
			golden.Update(t, testdata, name, tt.want)
			diff := golden.Diff(t, testdata, name, tt.got)
			switch {
			case diff == "" && tt.hasDiff:
				t.Error("there are any expected differences")
			case diff != "" && !tt.hasDiff:
				t.Error("there are some unexpected differences:", diff)
			}
		})
	}
}
