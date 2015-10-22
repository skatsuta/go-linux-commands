package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const usage = "usage: %v path\n"

func main() {
	if len(os.Args) != 2 {
		printErr(usage, os.Args[0])
	}

	if e := dirname(os.Stdout, os.Args[1]); e != nil {
		printErr(e.Error())
	}
}

// dirname prints a directory name of path to out.
func dirname(out io.Writer, path string) error {
	dir := filepath.Dir(path)
	_, err := fmt.Fprintln(out, dir)
	return err
}

// printErr prints a message to stderr and exit 1.
func printErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
