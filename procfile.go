package main

import (
	"math"
	"os"
	"regexp"

	_ "github.com/ddollar/forego/Godeps/_workspace/src/github.com/kr/pretty"
)

var procfileEntryRegexp = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")

type ProcfileEntry struct {
	Name    string
	Command string
}

type Procfile struct {
	Entries []ProcfileEntry
}

func ReadProcfile(filename string) (*Procfile, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return parseProcfile(fd)
}

func (pf *Procfile) HasProcess(name string) (exists bool) {
	for _, entry := range pf.Entries {
		if name == entry.Name {
			return true
		}
	}
	return false
}

func (pf *Procfile) LongestProcessName(concurrency map[string]int) (longest int) {
	longest = 6 // length of forego
	for _, entry := range pf.Entries {
		thisLen := len(entry.Name)
		// The "."
		thisLen += 1
		if c, ok := concurrency[entry.Name]; ok {
			// Add the number of digits
			thisLen += int(math.Log10(float64(c))) + 1
		}
		if thisLen > longest {
			longest = thisLen
		}
	}
	return
}
