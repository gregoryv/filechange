package filechange

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gregoryv/working"
)

func TestWatch_long(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	d := new(working.Directory)
	d.Temporary()
	defer d.RemoveAll()

	var (
		calls    int
		multiple bool
		sens     = new(Sensor)
	)
	sens.Root = d.Path()
	sens.Pause = 100 * time.Millisecond
	sens.Visit = func(modified ...string) {
		calls++
		if len(modified) > 1 {
			multiple = true
		}
		d.Touch("y")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go sens.Run(ctx)
	time.Sleep(100 * time.Millisecond)
	d.Touch("x")

	plus := 3 * sens.Pause
	time.Sleep(plus)
	d.Touch("x")
	time.Sleep(plus)
	if calls != 2 {
		t.Errorf("File changed twice but sensor reacted %v times", calls)
	}
	if multiple {
		t.Error("Got multiple changes")
	}
}

func TestWatch(t *testing.T) {
	d := new(working.Directory)
	d.Temporary()
	defer d.RemoveAll()
	d.MkdirAll("sub", "vendor/a/b")
	var (
		called bool
		sens   = new(Sensor)
	)
	sens.UseDefaults()
	sens.Root = d.Path()
	sens.Visit = func(...string) { called = true }
	sens.Pause = 50 * time.Millisecond
	plus := sens.Pause + 10*time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	go sens.Run(ctx)
	defer cancel()
	time.Sleep(plus)

	shouldSense := func(s string, err error) {
		t.Helper()
		called = false
		time.Sleep(plus)
		if !called {
			t.Error(s)
		}
	}
	shouldSense(d.Touch("a"))

	shouldNotSense := func(s string, err error) {
		t.Helper()
		called = false
		time.Sleep(plus)
		if called {
			t.Error(s)
		}
	}
	// Not recursive
	shouldNotSense(d.Touch("sub/hello"))

	// vendor should be ignored by default
	shouldNotSense(d.Touch("vendor/noop"))

	// Directories are ignored by default
	shouldNotSense(d.Touch("vendor"))

	// Removed
	sens.Recursive = true
	d.MkdirAll("sub")
	d.Touch("sub/x")
	os.RemoveAll(d.Join("sub"))
	time.Sleep(plus)
	shouldNotSense("", nil)
}
