package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
)

var procfileEntryRegexp = regexp.MustCompile("^([A-Za-z0-9_-]+):\\s*(.+)$")

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
		} else {
			// The index number after the dot.
			thisLen += 1
		}
		if thisLen > longest {
			longest = thisLen
		}
	}
	return
}

func parseProcfile(r io.Reader) (*Procfile, error) {
	pf := new(Procfile)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := procfileEntryRegexp.FindStringSubmatch(scanner.Text())
		if len(parts) > 0 {
			pf.Entries = append(pf.Entries, ProcfileEntry{parts[1], parts[2]})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Reading Procfile: %s", err)
	}
	return pf, nil
}
