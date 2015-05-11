package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const usage = "usage: %v path\n"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		return
	}

	dir := filepath.Dir(os.Args[1])
	fmt.Println(dir)
}
