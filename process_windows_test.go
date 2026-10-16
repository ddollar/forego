package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessQuotedCommand(t *testing.T) {
	process := NewProcess(t.TempDir(), `echo "owned fixture"`, Env{}, false)
	stdout, err := process.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := process.Start(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	if err := process.Wait(); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(output)); got != `"owned fixture"` {
		t.Fatalf("output = %q, want quoted fixture", got)
	}
}

func TestProcessWorkingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "working directory")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	process := NewProcess(dir, "cd", Env{}, false)
	stdout, err := process.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := process.Start(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	if err := process.Wait(); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(output)); !strings.EqualFold(filepath.Clean(got), filepath.Clean(dir)) {
		t.Fatalf("working directory = %q, want %q", got, dir)
	}
}
