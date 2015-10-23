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

	// read stdin
	if l == 0 {
		err := head(os.Stdin, os.Stdout, nlines)
		if err != nil {
			printErr(err)
			return
		}
	}

	// read files
	for _, path := range flag.Args() {
		fp, err := os.Open(path)
		if err != nil {
			printErr(err)
			return
		}
		defer func() {
			if e := fp.Close(); e != nil {
				printErr(e)
			}
		}()

		if e := head(fp, os.Stdout, nlines); e != nil {
			printErr(e)
			return
		}
	}
}

func head(in io.Reader, out io.Writer, n int) error {
	if n <= 0 {
		return fmt.Errorf("illegal line count -- %d", n)
	}

	r := bufio.NewReader(in)
	w := bufio.NewWriter(out)

	for {
		c, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		if e := w.WriteByte(c); e != nil {
			return e
		}

		if c == '\n' {
			if e := w.Flush(); e != nil {
				return e
			}

			n--
			if n == 0 {
				return nil
			}
		}
	}

	return w.Flush()
}

func printErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "head: %v\n", err)
}
