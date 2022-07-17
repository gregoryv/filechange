package filechange

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Sensor struct {
	Recursive bool
	Pause     time.Duration // between scans
	Last      time.Time
	Ignore    []string // filenames to ignore
	Root      string
	Visit     Visitor

	modified []string
}

// Run blocks until context is done and should only be called once.
func (s *Sensor) Run(ctx context.Context) {
	if s.Pause == 0 { // make sure we don't spin out of control
		s.Pause = time.Second
	}
	s.Last = time.Now()
	s.modified = make([]string, 0)
	if s.Visit == nil {
		s.Visit = noop
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.Pause):
			s.scanForChanges()
			if len(s.modified) > 0 {
				s.Visit(s.modified...)
				// Reset modified files, should not leak memory as
				// it's only strings
				s.modified = s.modified[:0:0]
				s.Last = time.Now()
			}
		}
	}
}

// Visitor is the function called with all modified filenames as
// arguments.
type Visitor func(modified ...string)

func noop(...string) {}

func (s *Sensor) scanForChanges() {
	if s.Root == "" {
		s.Root, _ = os.Getwd()
	}
	filepath.Walk(s.Root, s.checkModTime)
}

// checkModTime checks the given path and file info if it was modified
// after the last check. All modified paths are stored in s.modified.
// Directory changes are ignored and configured s.Ignore paths.
func (s *Sensor) checkModTime(path string, f os.FileInfo, err error) error {
	if f == nil {
		return nil
	}
	if s.ignore(path, f) {
		if f.IsDir() {
			if s.Recursive {
				return nil
			}
			return filepath.SkipDir
		}
		return nil
	}
	if !f.IsDir() && f.ModTime().After(s.Last) {
		s.modified = append(s.modified, path)
	}
	if !f.IsDir() {
		return nil
	}
	if s.Root == path {
		// the starting directory
		return nil
	}
	if s.Recursive {
		return nil
	}
	return filepath.SkipDir
}

// Ignore returns true if the file should be ignored
func (s *Sensor) ignore(path string, f os.FileInfo) bool {
	for _, str := range s.Ignore {
		if strings.Contains(path, str) {
			return true
		}
	}
	return false
}
