package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// flags
	notCreate := flag.Bool("c", false, "do not create the file if it does not exist")
	flag.Parse()

	args := flag.Args()

	// check if files exist
	for _, f := range args {
		if _, e := os.Lstat(f); os.IsNotExist(e) {
			if !*notCreate {
				if _, err := os.Create(f); err != nil {
					perror(err)
				}
			}
			continue
		}

		now := time.Now()
		if e := os.Chtimes(f, now, now); e != nil {
			perror(e)
		}
	}
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
}
