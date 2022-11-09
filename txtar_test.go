package golden_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tenntenn/golden"
)

func TestDirInit(t *testing.T) {
	t.Parallel()
	__ := func(args ...string) string {
		return golden.TxtarWith(t, args...)
	}
	cases := []struct {
		name    string
		initstr string
	}{
		{"single", __("a.txt", "hello")},
		{"multi", __("a.txt", "hello", "b.txt", "hi")},
		{"directory", __("dir/a.txt", "hello", "b.txt", "hi")},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			golden.DirInit(t, dir, tt.initstr)
			got := golden.Txtar(t, dir)
			if diff := cmp.Diff(tt.initstr, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestTxtarJoin(t *testing.T) {
	t.Parallel()
	__ := func(args ...string) string {
		return golden.TxtarWith(t, args...)
	}

	cases := map[string][]string{
		"empty":     {},
		"single":    {__("a.txt", "hello")},
		"multi":     {__("a.txt", "hello"), __("b.txt", "hi")},
		"directory": {__("dir/a.txt", "hello"), __("dir/b.txt", "hi")},
		"same":      {__("a.txt", "hello"), __("a.txt", "HELLO")},
		"comment":   {"comment-b\n" + __("b.txt", "hi"), "comment-a\n" + __("a.txt", "hello")},
	}

	testdata := filepath.Join("testdata", t.Name())
	for name, txtars := range cases {
		name, txtars := name, txtars
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := golden.TxtarJoin(t, txtars...)
			if diff := golden.Check(t, flagUpdate, testdata, name, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
