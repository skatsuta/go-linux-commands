package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	dispEnd = `$\n`
	dispTab = "> "
)

type replacer interface {
	replace(b byte) string
}

type disp struct {
	end bool
	tab bool
}

func (d *disp) replace(b byte) string {
	if b == '\n' && d.end {
		return dispEnd
	}

	if b == '\t' && d.tab {
		return dispTab
	}

	return string(b)
}

func newDisp() *disp {
	d := new(disp)

	// define flags
	flag.BoolVar(&d.end, "e", false, "display a dollar sign (`"+dispEnd+"`) at the end of each line")
	flag.BoolVar(&d.tab, "t", false, "display tab characters as `"+dispTab+"`")
	flag.Parse()
	flag.Usage = usage

	return d
}

func main() {
	d := newDisp()

	files := flag.Args()

	if len(files) < 1 {
		if e := doCat(os.Stdin, d); e != nil {
			die(e)
		}
		return
	}

	for _, file := range files {
		fp, err := os.Open(file)
		if err != nil {
			die(err)
		}
		defer func() {
			if e := fp.Close(); e != nil {
				die(e)
			}
		}()

		if e := doCat(fp, d); e != nil {
			die(e)
		}
	}
}

func doCat(fp *os.File, rpl replacer) error {
	r := bufio.NewReader(fp)
	w := bufio.NewWriter(os.Stdout)

	for {
		c, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		s := rpl.replace(c)

		if _, e := w.WriteString(s); e != nil {
			return e
		}

		if c == '\n' {
			if e := w.Flush(); e != nil {
				return e
			}
		}
	}

	return w.Flush()
}

func die(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s [-et] [file ...]:\n", os.Args[0])
	flag.PrintDefaults()
}
