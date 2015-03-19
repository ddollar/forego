package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

func parseProcfile(r io.Reader) (*Procfile, error) {
	pf := new(Procfile)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := envNixToWin(scanner.Bytes())
		parts := procfileEntryRegexp.FindStringSubmatch(line)
		if len(parts) > 0 {
			pf.Entries = append(pf.Entries, ProcfileEntry{parts[1], parts[2]})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Reading Procfile: %s", err)
	}
	return pf, nil
}

func envNixToWin(line []byte) string {
	nr := regexp.MustCompile("\\${?(\\w+)}?")
	out := nr.ReplaceAll(line, []byte("%$1%"))
	return string(out)
}
