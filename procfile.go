package main

import (
  "bytes"
  "github.com/kr/pretty"
  "io"
  "io/ioutil"
  "os"
  "regexp"
)

const (
  procfileEntryRegexp = "^([A-Za-z0-9_]+):\\s*(.+)$"
)

type ProcfileEntry struct {
  name string
  command string
}

type Procfile struct {
  entries []ProcfileEntry
}

var _ = pretty.Println // lol

func OpenProcfile(filename string) (*Procfile, error) {
  fd, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer fd.Close()
  return parseProcfile(fd)
}

func parseProcfile(r io.Reader) (*Procfile, error) {
  pf := new(Procfile)
  b, err := ioutil.ReadAll(r)
  if err != nil {
    return nil, err
  }
  lines := bytes.Split(b, []byte("\n"))
  for _, line := range lines {
    r := regexp.MustCompile(procfileEntryRegexp)
    parts := r.FindStringSubmatch(string(line))
    if parts != nil {
      pf.entries = append(pf.entries, ProcfileEntry{parts[1], parts[2]})
    }
  }
  return pf, nil
}
