package main

import (
	"fmt"
	"os"
)

func main() {
	s := "y"

	if len(os.Args) > 1 {
		s = os.Args[1]
	}

	for {
		fmt.Println(s)
	}
}
