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
	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

// Diff compares between the given data and a golden file which is stored in testdata as name+".golden".
// Diff returns difference of them.
func Diff(t *testing.T, testdata, name string, data interface{}) string {
	t.Helper()
	path := filepath.Join(testdata, name+".golden")
	golden, err := os.Open(path)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer golden.Close()

	want, got := readAll(t, golden), readAll(t, data)

	return cmp.Diff(want, got)
}

func readAll(t *testing.T, data interface{}) string {
	r := newReader(t, data)
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	return string(b)
}

func newReader(t *testing.T, data interface{}) io.Reader {
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
func Update(t *testing.T, testdata, name string, data interface{}) {
	t.Helper()
	path := filepath.Join(testdata, name+".golden")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	r := newReader(t, data)
	if _, err := io.Copy(f, r); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := f.Close(); err != nil {
		t.Fatal("unexpected error:", err)
	}
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

// Txtar converts a directory as a txtar format.
func Txtar(t *testing.T, dir string) string {
	t.Helper()
	ar, err := txtarfs.From(os.DirFS(dir))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	return string(txtar.Format(ar))
}
