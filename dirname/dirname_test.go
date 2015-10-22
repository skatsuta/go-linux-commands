package main

import (
	"bytes"
	"testing"
)

func TestDirname(t *testing.T) {
	testCases := []struct {
		path string
		want string
	}{
		{"", "."},
		{".", "."},
		{"./foo", "."},
		{"foo", "."},
		{"foo/bar", "foo"},
		{"/foo/bar", "/foo"},
		{"foo/bar/baz", "foo/bar"},
		{"/foo/bar/baz", "/foo/bar"},
	}

	var out bytes.Buffer
	for _, tt := range testCases {
		if e := dirname(&out, tt.path); e != nil {
			t.Fatalf("dirname failed: %v", e)
		}

		got := out.String()
		if got != tt.want+"\n" {
			t.Errorf("path = %q: got = %q; want = %q", tt.path, got, tt.want)
		}

		out.Reset()
	}
}
