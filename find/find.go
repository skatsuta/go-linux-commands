package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	doFind(os.Stdout, os.Stderr, flag.Args()...)
}

func doFind(stdout, stderr io.Writer, paths ...string) {
	for _, path := range paths {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintln(stderr, err)
				return nil
			}

			fmt.Fprintln(stdout, p)
			return nil
		})

		if err != nil {
			fmt.Fprintln(stderr, err)
		}
	}
}
