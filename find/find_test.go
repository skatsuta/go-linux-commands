package main

import (
	"bytes"
	"testing"
)

func TestFind(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	paths := []string{"a", "find.go"}
	doFind(stdout, stderr, paths...)

	wantOut := "find.go\n"
	wantErr := "lstat a: no such file or directory\n"

	if stdout.String() != wantOut {
		t.Errorf("got: %s; want: %s", stdout.String(), wantOut)
	}

	if stderr.String() != wantErr {
		t.Errorf("got: %s; want: %s", stderr.String(), wantErr)
	}
}
