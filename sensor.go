// Package filechange provides a sensor of file modifications.
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
		Interval: 2 * time.Second,
		root:     root,
		visit:    v,
	}
}

type Sensor struct {
	Recursive bool
	// between scans, should be > 0, sane values are >1s
	Interval time.Duration
	Last     time.Time
	Ignore   []string // filenames to ignore

	visit    Visitor
	root     string
	modified []string
}

// Run blocks until context is done and should only be called once.
func (s *Sensor) Run(ctx context.Context) {
	s.Last = time.Now()
	s.modified = make([]string, 0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.Interval):
			s.scanForChanges()
			if len(s.modified) > 0 {
				s.visit(s.modified...)
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
