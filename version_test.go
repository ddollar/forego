package main

import (
	"bytes"
	"testing"
)

func TestVersion(t *testing.T) {
	var b bytes.Buffer
	stdout = &b
	cmdVersion.Run(cmdVersion, []string{})
	output := b.String()
	assertEqual(t, output, "dev\n")
}

func assertEqual(t *testing.T, a, b interface{}) {
	if a != b {
		t.Fatalf(`Expected %#v to equal %#v`, a, b)
	}
}
