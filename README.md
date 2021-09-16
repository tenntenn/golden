# golden

[![pkg.go.dev][gopkg-badge]][gopkg]

`golden` provides utilities for golden file tests.

```go
package a_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/tenntenn/golden"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func testTarget(dir string) error {
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0700); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world"), 0700); err != nil {
		return err
	}

	return nil
}

func Test(t *testing.T) {
	dir := t.TempDir()
	if err := testTarget(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(dir)

	if flagUpdate {
		golden.Update(t, "testdata", "mytest", got)
		return
	}

	if diff := golden.Diff(t, "testdata", "mytest", got); diff != "" {
		t.Error(diff)
	}
}
```

<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/tenntenn/golden
[gopkg-badge]: https://pkg.go.dev/badge/github.com/tenntenn/golden?status.svg
