// +build !windows

package main

import (
	"bufio"
	"fmt"
	"io"
)

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
