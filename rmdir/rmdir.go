package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		perror(fmt.Errorf("%s: no arguments\n", os.Args[0]))
	}

	for _, path := range os.Args {
		e := os.Remove(path)
		if e != nil {
			perror(e)
		}
	}
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "[error] %v\n", err)
}
