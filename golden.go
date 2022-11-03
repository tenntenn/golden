package golden

import (
	"bytes"
	"encoding"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Check updates a golden file when update is true otherwise compares data with the exsiting golden file by DiffWithOpts.
// If update is true Check does not compare and just return "".
//
//	var flagUpdate bool
//
//	func init() {
//		flag.BoolVar(&flagUpdate, "update", false, "update golden files")
//	}
//
//	func Test(t *testing.T) {
//		got := doSomething()
//		if diff := golden.Check(t, flagUpdate, "testdata", t.Name(), got); diff != "" {
//			t.Error(diff)
//		}
//	}
func Check(t *testing.T, update bool, testdata, name string, data any, opts ...cmp.Option) string {
	t.Helper()
	if update {
		Update(t, testdata, name, data)
		return ""
	}
	return DiffWithOpts(t, testdata, name, data, opts...)
}

// DiffWithOpts compares between the given data and a golden file which is stored in testdata as name+".golden".
// DiffWithOpts returns difference of them.
// DiffWithOpts uses [go-cmp] to compare.
//
// [go-cmp]: https://pkg.go.dev/github.com/google/go-cmp/cmp
func DiffWithOpts(t *testing.T, testdata, name string, data any, opts ...cmp.Option) string {
	t.Helper()
	return New(t, false, testdata, name, opts...).Check("", data)
}

// Diff compares between the given data and a golden file which is stored in testdata as name+".golden".
// Diff returns difference of them.
// Diff uses [go-cmp] to compare.
//
// [go-cmp]: https://pkg.go.dev/github.com/google/go-cmp/cmp
func Diff(t *testing.T, testdata, name string, data any) string {
	t.Helper()
	return DiffWithOpts(t, testdata, name, data, nil)
}

func readAll(t TestingT, data any) string {
	t.Helper()
	r := newReader(t, data)
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	return string(b)
}

func newReader(t TestingT, data any) io.Reader {
	t.Helper()
	switch data := data.(type) {
	case io.Reader:
		return data
	case string:
		return strings.NewReader(data)
	case []byte:
		return bytes.NewReader(data)
	case encoding.TextMarshaler:
		b, err := data.MarshalText()
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		return bytes.NewReader(b)
	default:
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			t.Fatal("unexpected error:", err)
		}
		return &buf
	}
}

// Update updates a golden file with the given data.
// The golden file saved as name+".golden" in testdata.
func Update(t *testing.T, testdata, name string, data any) {
	t.Helper()
	_ = New(t, true, testdata, name).Check("", data)
}

// RemoveAll removes all golden files which has .golden extention and is under testdata.
func RemoveAll(t *testing.T, testdata string) {
	t.Helper()

	err := filepath.Walk(testdata, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".golden" {
			return nil
		}

		if err := os.Remove(path); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Fatal("unexpected error", err)
	}
}
