package filechange

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestSensor_Run(t *testing.T) {
	// test changes in a temporary directory
	dir := t.TempDir()

	// create a structure
	os.MkdirAll(filepath.Join(dir, "sub"), 0722)
	os.MkdirAll(filepath.Join(dir, "vendor/a/b"), 0722)

	var (
		called bool
		sens   = &Sensor{
			Ignore:    []string{"vendor/"},
			Recursive: true,
			Root:      dir,
			Visit:     func(...string) { called = true },
			Pause:     50 * time.Millisecond,
		}
		plus  = sens.Pause + 10*time.Millisecond
		touch = func(filename string) *exec.Cmd {
			cmd := exec.Command("touch", filepath.Join(dir, filename))
			if err := cmd.Run(); err != nil {
				t.Helper()
				t.Fatal(err)
			}
			return cmd
		}
	)

	// start sensor
	ctx, cancel := context.WithCancel(context.Background())
	go sens.Run(ctx)
	defer cancel()
	time.Sleep(plus)

	// create file in root triggers sensor
	cmd := touch("a.txt")
	if time.Sleep(plus); !called {
		t.Errorf("%q should trigger sensor", cmd)
	}

	// create file in subdir triggers sensor when recursive
	called = false // reset
	cmd = touch("sub/b.txt")
	if time.Sleep(plus); !called {
		t.Errorf("%q should trigger sensor", cmd)
	}

	// create file in ignored subdir does not trigger sensor
	called = false // reset
	cmd = touch("vendor/noop")
	if time.Sleep(plus); called {
		t.Errorf("%q triggered sensor on ignored directory", cmd)
	}

	// create directory is ignored
	called = false
	os.MkdirAll(filepath.Join(dir, "Xdir"), 0722)
	if time.Sleep(plus); called {
		t.Errorf("mkdir in root triggered sensor")
	}

	/*

		// Removed
		sens.Recursive = true
		d.MkdirAll("sub")
		d.Touch("sub/x")
		os.RemoveAll(d.Join("sub"))
		time.Sleep(plus)
		shouldNotSense("", nil)
	*/
}
