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
	// test sensor in a temporary directory
	dir := t.TempDir()

	// create a structure
	os.MkdirAll(filepath.Join(dir, "sub"), 0722)
	os.MkdirAll(filepath.Join(dir, "vendor/a/b"), 0722)

	var called bool
	sensor := NewSensor(dir, func(...string) { called = true })
	sensor.Ignore = []string{"vendor/", "build"}
	sensor.Recursive = true
	// Use shorter interval to speed up test
	walkEstimate := time.Millisecond // measured to 29.098µs
	sensor.Interval = 10 * walkEstimate

	var (
		plus  = sensor.Interval + 3*walkEstimate
		touch = func(filename string) *exec.Cmd {
			cmd := exec.Command("touch", filepath.Join(dir, filename))
			if err := cmd.Run(); err != nil {
				t.Helper()
				t.Fatal(err)
			}
			return cmd
		}
		reset = func() { called = false }
		wait  = func() bool { time.Sleep(plus); return true }
	)

	// Note! the same sensor instance is used in all cases
	// start sensor
	ctx, cancel := context.WithCancel(context.Background())
	go sensor.Run(ctx)
	defer cancel()
	wait()

	// create file in root triggers sensor
	reset()
	if cmd := touch("a.txt"); wait() && !called {
		t.Errorf("%q should trigger sensor", cmd)
	}

	// "create file in subdir triggers sensor when recursive"
	reset()
	if cmd := touch("sub/b.txt"); wait() && !called {
		t.Errorf("%q should trigger sensor", cmd)
	}

	// create file in ignored subdir does not trigger sensor
	reset()
	cmd := touch("vendor/noop")
	if wait() && called {
		t.Errorf("%q triggered sensor on ignored directory", cmd)
	}

	// create directory is ignored
	reset()
	os.MkdirAll(filepath.Join(dir, "Xdir"), 0722)
	if wait() && called {
		t.Errorf("mkdir in root triggered sensor")
	}

	// create file in subdir is ignored when Not recursive
	reset()
	sensor.Recursive = false

	if cmd := touch("sub/something.txt"); wait() && called {
		t.Errorf("%q triggered sensor", cmd)
	}

	// create directory in root is ignored
	reset()
	sensor.Recursive = false
	os.MkdirAll(filepath.Join(dir, "build"), 0722)
	if wait() && called {
		t.Errorf("%q triggered sensor", cmd)
	}

	// FileInfo is nil
	reset()
	sensor.root = "/no-such-directory"
	if wait() && called {
		t.Errorf("%q triggered sensor", cmd)
	}

	// Removing file should not trigger sensor
	reset()
	sensor.Recursive = true
	_ = touch("sub/toremove")
	os.RemoveAll(filepath.Join(dir, "sub/toremove"))
	if wait() && called {
		t.Error("removing file triggered sensor")
	}
} // 106
