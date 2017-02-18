package main

import (
	"bytes"
	"strings"
	"testing"
)

var testSrcs = []string{
	`
	line1	
	line2`,
	`
	line3	
	line4
	`,
}

func TestCat(t *testing.T) {
	want := `
	line1	
	line2
	line3	
	line4
	`

	var out bytes.Buffer

	for _, src := range testSrcs {
		r := strings.NewReader(src)
		if e := cat(r, &out); e != nil {
			t.Errorf("cat error: %v", e)
		}
	}

	if got := out.String(); want != got {
		t.Errorf("want: %s\ngot: %s", want, got)
	}
}

func TestCatShowAll(t *testing.T) {
	// Show all
	*showEnds = true
	*showTabs = true

	want := `$
^Iline1^I$
^Iline2$
^Iline3^I$
^Iline4$
^I`

	var out bytes.Buffer

	for _, src := range testSrcs {
		r := strings.NewReader(src)
		if e := cat(r, &out); e != nil {
			t.Errorf("cat error: %v", e)
		}
	}

	if got := out.String(); want != got {
		t.Errorf("want: %s\ngot: %s", want, got)
	}
}
