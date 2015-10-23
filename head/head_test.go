package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHead(t *testing.T) {
	testCases := []struct {
		in   string
		n    int
		want string
	}{
		{"", 1, ""},
		{"foo", 1, "foo"},
		{"foo\nbar", 1, "foo\n"},
		{"foo\nbar", 2, "foo\nbar"},
	}

	var out bytes.Buffer
	for _, tt := range testCases {
		in := strings.NewReader(tt.in)
		if e := head(in, &out, tt.n); e != nil {
			t.Errorf("head failed: %v", e)
		}

		got := out.String()
		if got != tt.want {
			t.Errorf("in = %q, n = %d: got = %q; want = %q", tt.in, tt.n, got, tt.want)
		}

		out.Reset()
	}
}

func TestHeadErr(t *testing.T) {
	testCases := []struct {
		in string
		n  int
	}{
		{"foo", -1},
		{"foo", 0},
	}

	var out bytes.Buffer
	for _, tt := range testCases {
		in := strings.NewReader(tt.in)
		if e := head(in, &out, tt.n); e == nil {
			t.Errorf("in = %q, n = %d: error should occur but got <nil>", tt.in, tt.n)
		}

		out.Reset()
	}
}
