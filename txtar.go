package golden

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

// DirInit creates directory by txtar format.
func DirInit(t *testing.T, root, txtarStr string) {
	t.Helper()
	fsys := txtarfs.As(txtar.Parse([]byte(txtarStr)))
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) (rerr error) {
		if err != nil {
			return err
		}

		// directory would create with a file
		if d.IsDir() {
			return nil
		}

		dstPath := filepath.Join(root, filepath.FromSlash(path))

		src, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			if err := src.Close(); err != nil && rerr == nil {
				rerr = err
			}
		}()

		fi, err := src.Stat()
		if err != nil {
			return err
		}

		if fi.Size() == 0 {
			return nil
		}

		err = os.MkdirAll(filepath.Dir(dstPath), 0700)
		if err != nil {
			return err
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer func() {
			if err := dst.Close(); err != nil && rerr == nil {
				rerr = err
			}
		}()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

// Txtar converts a directory as a txtar format.
func Txtar(t *testing.T, dir string) string {
	t.Helper()
	ar, err := txtarfs.From(os.DirFS(dir))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	sortTxtar(ar)
	return string(txtar.Format(ar))
}

// TxtarWith creates a txtar format value with given a file name and its data pairs.
func TxtarWith(t *testing.T, nameAndData ...string) string {
	t.Helper()
	if len(nameAndData)%2 != 0 {
		t.Fatal("invalid argument:", nameAndData)
	}

	ar := &txtar.Archive{
		Files: make([]txtar.File, 0, len(nameAndData)/2),
	}

	for i := 0; i < len(nameAndData); i += 2 {
		ar.Files = append(ar.Files, txtar.File{
			Name: nameAndData[i],
			Data: []byte(nameAndData[i+1]),
		})
	}
	sortTxtar(ar)

	return string(txtar.Format(ar))
}

func sortTxtar(ar *txtar.Archive) {
	sort.Slice(ar.Files, func(i, j int) bool {
		return ar.Files[i].Name < ar.Files[j].Name
	})
}
