package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
)

func TestLineReaderFinalOutput(t *testing.T) {
	for _, stderr := range []bool{false, true} {
		for _, input := range []string{"", "complete\n", "tail", "complete\ntail", strings.Repeat("x", 2048)} {
			for _, withEOF := range []bool{false, true} {
				var reader io.Reader = strings.NewReader(input)
				if withEOF {
					reader = iotest.DataErrReader(reader)
				}
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatal(err)
				}
				stdout := os.Stdout
				os.Stdout = w
				var wg sync.WaitGroup
				wg.Add(1)
				NewOutletFactory().LineReader(&wg, "web.1", 0, reader, stderr)
				wg.Wait()
				os.Stdout = stdout
				w.Close()
				output, err := io.ReadAll(r)
				r.Close()
				if err != nil {
					t.Fatal(err)
				}
				lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
				if input == "" {
					lines = nil
				}
				if got := strings.Count(string(output), "web.1"); got != len(lines) {
					t.Fatalf("input %q stderr=%t withEOF=%t: got %d lines, want %d: %q", input, stderr, withEOF, got, len(lines), output)
				}
				for _, line := range lines {
					if !strings.Contains(string(output), line) {
						t.Fatalf("input %q stderr=%t withEOF=%t: missing %q in %q", input, stderr, withEOF, line, output)
					}
				}
			}
		}
	}
}
