package golden

import (
	"io"
	"os"
	"path/filepath"

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
	testingT TestingT
	update   bool
	testdata string
	name     string
	opts     []cmp.Option
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

	want, got := readAll(c.testingT, golden), readAll(c.testingT, data)

	return cmp.Diff(want, got, c.opts...)
}

func (c *Checker) updateFile(path string, data any) {
	c.testingT.Helper()

	f, err := os.Create(path)
	if err != nil {
		c.testingT.Fatal("unexpected error:", err)
	}

	r := newReader(c.testingT, data)
	if _, err := io.Copy(f, r); err != nil {
		c.testingT.Fatal("unexpected error:", err)
	}

	if err := f.Close(); err != nil {
		c.testingT.Fatal("unexpected error:", err)
	}
}

// CheckAll do golden file tests for multiple data.
// CheckAll calls [Checker.Check] for each data and the handler function with the return value of Check.
// The handler should not call t.Helper because CheckAll calls it.
//
//	var flagUpdate bool
//
//	func init() {
//		flag.BoolVar(&flagUpdate, "update", false, "update golden files")
//	}
//
//	func Test(t *testing.T) {
//		var got1 bytes.Buffer
//		got2 := doSomething(&got1)
//		c := golden.New(t, flagUpdate, "testdata", t.Name())
//		c.CheckAll(map[string]any{
//			"_got1": &got1,
//			"_got2": got2,
//		}, func(diff string) {
//			if diff != "" {
//				t.Error(diff)
//			}
//		})
//	}
func (c *Checker) CheckAll(suffixToData map[string]any, handler func(diff string)) {
	c.testingT.Helper()
	for suffix, data := range suffixToData {
		if diff := c.Check(suffix, data); diff != "" {
			c.testingT.Helper()
			handler(diff)
		}
	}
}
