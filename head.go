package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	defaultNumLines = 10
)

func main() {
	// Define usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [-n LINES] [file ...]:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// process flags
	nlines := defaultNumLines
	flag.IntVar(&nlines, "n", defaultNumLines, "the number of lines to display")
	flag.Parse()

	l := len(flag.Args())

	if l == 0 {
		doHead(os.Stdin, nlines)
		return
	}

	for _, path := range flag.Args() {
		fp, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}
		defer func() {
			if e := fp.Close(); e != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				return
			}
		}()

		doHead(fp, nlines)
	}
}

func doHead(fp *os.File, nlines int) {
	if nlines <= 0 {
		return
	}

	r := bufio.NewReader(fp)
	w := bufio.NewWriter(os.Stdout)

	for {
		c, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		if e := w.WriteByte(c); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			os.Exit(1)
		}

		if c == '\n' {
			if e := w.Flush(); e != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				return
			}

			nlines--
			if nlines == 0 {
				return
			}
		}
	}

	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
}
