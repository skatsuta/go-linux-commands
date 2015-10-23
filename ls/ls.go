package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	unixHiddenFilePrefix = "."
	currentDir           = "."
)

type option struct {
	all bool
	rec bool
}

func newOption() option {
	opt := option{}

	flag.BoolVar(&opt.all, "A", false, "list all entries except for . and ..")
	flag.BoolVar(&opt.rec, "R", false, "recursively list subdirectories encountered")

	flag.Parse()

	return opt
}

func main() {
	opt := newOption()

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{currentDir}
	}

	for _, path := range paths {
		if e := ls(os.Stdout, path, opt); e != nil {
			perror(e)
		}
	}
}

func ls(out io.Writer, path string, opt option) error {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		name := fi.Name()

		// skip the following process if not to display hidden files
		skip := !opt.all && isHidden(name)
		if skip {
			continue
		}

		fmt.Fprintln(out, name)

		// run recursively
		isRec := opt.rec && fi.IsDir()
		if isRec {
			name = filepath.Join(path, name)
			if e := ls(out, name, opt); e != nil {
				return e
			}
		}
	}

	return nil
}

func isHidden(path string) bool {
	return strings.HasPrefix(path, unixHiddenFilePrefix)
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "ls: %v\n", err)
}
