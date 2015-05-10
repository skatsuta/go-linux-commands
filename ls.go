package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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
		doLs(path, opt)
	}
}

func doLs(path string, opt option) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		perror(err)
		return
	}

	for _, fi := range fis {
		name := fi.Name()

		// skip the following process if not to display hidden files
		skip := !opt.all && isHidden(name)
		if skip {
			continue
		}

		fmt.Println(name)

		// run recursively
		isRec := opt.rec && fi.IsDir()
		if isRec {
			name = path + "/" + name
			doLs(name, opt)
		}
	}
}

func isHidden(path string) bool {
	return strings.HasPrefix(path, unixHiddenFilePrefix)
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "[error] %v\n", err)
}
