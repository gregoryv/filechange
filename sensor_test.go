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
	visit := func(...string) { called = true }

	sens := NewSensor(dir, visit)
	sens.Ignore = []string{"vendor/", "build"}
	sens.Recursive = true
	// Use shorter interval to speed up test
	walkEstimate := time.Millisecond // measured to 29.098µs
	sens.Pause = 10 * walkEstimate

	var (
		plus  = sens.Pause + 3*walkEstimate
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

	// create file in subdir is ignored when Not recursive
	called = false // reset
	sens.Recursive = false
	cmd = touch("sub/something.txt")
	if time.Sleep(plus); called {
		t.Errorf("%q triggered sensor", cmd)
	}

	// create directory in root is ignored
	called = false // reset
	sens.Recursive = false
	os.MkdirAll(filepath.Join(dir, "build"), 0722)
	if time.Sleep(plus); called {
		t.Errorf("%q triggered sensor", cmd)
	}

	// FileInfo is nil
	called = false // reset
	sens.root = "/no-such-directory"
	if time.Sleep(plus); called {
		t.Errorf("%q triggered sensor", cmd)
	}

	// Removing file should not trigger sensor
	sens.Recursive = true
	touch("sub/toremove")
	os.RemoveAll(filepath.Join(dir, "sub/toremove"))
	if time.Sleep(plus); called {
		t.Error("removing file triggered sensor")
	}
}
