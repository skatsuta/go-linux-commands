package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		perror(fmt.Errorf("%s: no arguments\n", os.Args[0]))
		return
	}

	for _, path := range os.Args {
		if e := os.Mkdir(path, 0777); e != nil {
			perror(e)
		}
	}
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "[error] %v\n", err)
}
