package golden_test

import (
	"strings"
	"testing"

	"github.com/tenntenn/golden"
)

type marshaler string

func (m marshaler) MarshalText() (text []byte, err error) {
	return []byte(m), nil
}

func TestDiff(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		want    string
		got     any
		hasDiff bool
	}{
		{"string-nodiff", "hello", "hello", false},
		{"bytes-nodiff", "hello", []byte("hello"), false},
		{"reader-nodiff", "hello", strings.NewReader("hello"), false},
		{"json-nodiff", "{\"S\":\"hello\"}\n", struct{ S string }{S: "hello"}, false},
		{"marshaler-nodiff", "hello", marshaler("hello"), false},

		{"string-diff", "Hello", "hello", true},
		{"bytes-diff", "Hello", []byte("hello"), true},
		{"reader-diff", "Hello", strings.NewReader("hello"), true},
		{"json-diff", "{\"S\":\"Hello\"}\n", struct{ S string }{S: "hello"}, true},
		{"marshaler-diff", "Hello", marshaler("hello"), true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testdata := t.TempDir()
			golden.Update(t, testdata, tt.name, tt.want)
			diff := golden.Diff(t, testdata, tt.name, tt.got)
			switch {
			case diff == "" && tt.hasDiff:
				t.Error("there are any expected differences")
			case diff != "" && !tt.hasDiff:
				t.Error("there are some unexpected differences:", diff)
			}
		})
	}
}
