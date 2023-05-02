package golden

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// TestingT is interface for *testing.T.
type TestingT interface {
	Helper()
	Fatal(args ...any)
}

// Checker can do golden file testing for multiple data.
// Checker holds *testing.T, update flag, testdata directory,
// test name and options for [go-cmp].
//
// [go-cmp]: https://pkg.go.dev/github.com/google/go-cmp/cmp
type Checker struct {
	testingT  TestingT
	update    bool
	testdata  string
	name      string
	opts      []cmp.Option
	JSONIdent bool
}

// New creates a [Checker].
func New(t TestingT, update bool, testdata, name string, opts ...cmp.Option) *Checker {
	return &Checker{
		testingT: t,
		update:   update,
		testdata: testdata,
		name:     name,
		opts:     opts,
	}
}

// Check do a golden file test for a single data.
// Check calls [Check] function with test name which combiend with suffix.
//
//	var flagUpdate bool
//
//	func init() {
//		flag.BoolVar(&flagUpdate, "update", false, "update golden files")
//	}
//
//	func Test(t *testing.T) {
//		got := doSomething()
//		c := golden.New(t, flagUpdate, "testdata", t.Name())
//		if diff := c.Check("_someting", got); diff != "" {
//			t.Error(diff)
//		}
//	}
func (c *Checker) Check(suffix string, data any) (diff string) {
	c.testingT.Helper()

	path := filepath.Join(c.testdata, c.name+suffix+".golden")

	if c.update {
		c.updateFile(path, data)
		return ""
	}

	golden, err := os.Open(path)
	if err != nil {
		c.testingT.Fatal("unexpected error:", err)
	}
	defer golden.Close()

	wantStr := readAll(c.testingT, c.JSONIdent, golden)

	got := reflect.ValueOf(data)
	if got.Kind() != reflect.Pointer || got.Elem().Kind() != reflect.Struct {
		gotStr := readAll(c.testingT, c.JSONIdent, data)
		return cmp.Diff(wantStr, gotStr, c.opts...)
	}

	want := reflect.New(got.Elem().Type())
	dec := json.NewDecoder(strings.NewReader(wantStr))
	if err := dec.Decode(want.Interface()); err != nil {
		gotStr := readAll(c.testingT, c.JSONIdent, data)
		return cmp.Diff(wantStr, gotStr, c.opts...)
	}

	return cmp.Diff(want.Interface(), data, c.opts...)
}

func (c *Checker) updateFile(path string, data any) {
	c.testingT.Helper()

	f, err := os.Create(path)
	if err != nil {
		c.testingT.Fatal("unexpected error:", err)
	}

	r := newReader(c.testingT, c.JSONIdent, data)
	if _, err := io.Copy(f, r); err != nil {
		c.testingT.Fatal("unexpected error:", err)
	}

	if err := f.Close(); err != nil {
		c.testingT.Fatal("unexpected error:", err)
	}
}
