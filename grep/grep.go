package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	var ins, inv bool
	flag.BoolVar(&ins, "i", false, "Perform case insensitive matching")
	flag.BoolVar(&inv, "v", false, "Selected lines are those NOT matching any of the specified patterns")
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "no pattern\n")
		return
	}

	expr := args[0]
	if ins {
		expr = "(?i)" + expr
	}

	re, err := regexp.Compile(expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if len(args) == 1 {
		doGrep(re, os.Stdin, inv)
		return
	}

	for _, f := range args[1:] {
		fp, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		defer fp.Close()

		if e := doGrep(re, fp, inv); e != nil {
			fmt.Fprintf(os.Stderr, "%v\n", e)
			return
		}
	}
}

func doGrep(re *regexp.Regexp, fp *os.File, inv bool) error {
	r := bufio.NewReader(fp)
	w := bufio.NewWriter(os.Stdout)

	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		disp := re.MatchString(line) == !inv
		if disp {
			if _, e := w.WriteString(line); e != nil {
				return e
			}

			if e := w.Flush(); e != nil {
				return e
			}
		}
	}

	return nil
}
