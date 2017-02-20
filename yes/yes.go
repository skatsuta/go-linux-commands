package main

import (
	"fmt"
	"os"
	"strings"
)

const defaultArg = "y"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = append(args, defaultArg)
	}

	s := strings.Join(args, " ")
	for {
		fmt.Println(s)
	}
}
