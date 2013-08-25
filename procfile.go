package main

import (
  "bufio"
  "fmt"
  "github.com/kr/pretty"
  "io"
  "os"
  "regexp"
)

var procfileEntryRegexp = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")

type ProcfileEntry struct {
  Name string
  Command string
}

type Procfile struct {
  entries []ProcfileEntry
}

var _ = pretty.Println // lol

func ReadProcfile(filename string) (*Procfile, error) {
  fd, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer fd.Close()
  return parseProcfile(fd)
}

func parseProcfile(r io.Reader) (*Procfile, error) {
  pf := new(Procfile)
  scanner := bufio.NewScanner(r)
  for scanner.Scan() {
    parts := procfileEntryRegexp.FindStringSubmatch(scanner.Text())
    pf.entries = append(pf.entries, ProcfileEntry{parts[1], parts[2]})
  }
  if err := scanner.Err(); err != nil {
    return nil, fmt.Errorf("Reading Procfile: %s", err)
  }
  return pf, nil
}
