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
	modified  []string
	ignore    []string
	Root      string
	Visit     Visitor
}

func (s *Sensor) Run(ctx context.Context) {
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

type Visitor func(modified ...string)

// UseDefaults sets sensible sensor values and ignores for current working directory.
// React is a noop.
func (s *Sensor) UseDefaults() {
	s.Pause = time.Second
	s.Last = time.Now()
	s.modified = make([]string, 0)
	s.ignore = []string{"#", ".git/", "vendor/"}
	s.Root, _ = os.Getwd()
	s.Visit = noop
}

func noop(...string) {}

func (s *Sensor) scanForChanges() {
	filepath.Walk(s.Root, s.visit)
}

// Ignore returns true if the file should be ignored
func (s *Sensor) Ignore(path string, f os.FileInfo) bool {
	for _, str := range s.ignore {
		if strings.Contains(path, str) {
			return true
		}
	}
	return false
}

func (s *Sensor) visit(path string, f os.FileInfo, err error) error {
	if s.Ignore(path, f) {
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
