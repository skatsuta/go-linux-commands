package main

import (
	"fmt"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		perror(err)
	}

	fmt.Println(wd)
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "[error] %v\n", err)
}
