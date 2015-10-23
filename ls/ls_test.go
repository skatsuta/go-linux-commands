package main

import (
	"bytes"
	"testing"
)

func TestList(t *testing.T) {
	testCases := []struct {
		path string
		opt  option
		want string
	}{
		{".", option{}, "ls.go\nls_test.go\n"},
	}

	var out bytes.Buffer
	for _, tt := range testCases {
		if e := ls(&out, tt.path, tt.opt); e != nil {
			t.Fatalf("List failed: %v", e)
		}

		got := out.String()
		if got != tt.want {
			t.Errorf("path = %q, opt = %+v: got = %q; want = %q", tt.path, tt.opt, got, tt.want)
		}

		out.Reset()
	}
}
