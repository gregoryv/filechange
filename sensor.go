//Package filechange provides a sensor of file modifications.
package filechange

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewSensor(root string, v Visitor) *Sensor {
	return &Sensor{
		Pause: 2 * time.Second,
		c:     make(chan []string, 1),
		root:  root,
		visit: v,
	}
}

type Sensor struct {
	Recursive bool
	// between scans, should be > 0, sane values are >1s
	Pause  time.Duration
	Last   time.Time
	Ignore []string // filenames to ignore

	visit    Visitor
	root     string
	modified []string

	// used to signal last modifications, it is cleared between scans
	// see method Modified()
	c chan []string
}

// Run blocks until context is done and should only be called once.
func (s *Sensor) Run(ctx context.Context) {
	s.Last = time.Now()
	s.modified = make([]string, 0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.Pause):
			s.scanForChanges()
			if len(s.modified) > 0 {
				if s.visit != nil {
					s.visit(s.modified...)
				}
				modifiedFiles := make([]string, len(s.modified))
				copy(modifiedFiles, s.modified)
				// clear any old non read modifications
				select {
				case <-s.c:
				default:
				}
				// add new modifications
				s.c <- modifiedFiles
				// Reset modified files, should not leak memory as
				// it's only strings
				s.modified = s.modified[:0:0]
				s.Last = time.Now()
			}
		}
	}
}

func (s *Sensor) Modified() <-chan []string {
	return s.c
}

// Visitor is the function called with all modified filenames as
// arguments.
type Visitor func(modified ...string)

func (s *Sensor) scanForChanges() {
	filepath.Walk(s.root, s.checkModTime)
}

// checkModTime checks the given path and file info if it was modified
// after the last check. All modified paths are stored in s.modified.
// Directory changes are ignored and configured s.Ignore paths.
func (s *Sensor) checkModTime(path string, f os.FileInfo, err error) error {
	if err != nil { // handle case when Sensor.root doesn't exist
		return err
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
	if s.root == path {
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
