package golden_test

import (
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
