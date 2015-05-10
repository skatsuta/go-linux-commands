package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 3 {
		perror(fmt.Errorf("Usage: %s <command> <args>\n", os.Args[0]))
		return
	}

	cmd := os.Args[1]
	arg := os.Args[2]

	ch := make(chan string)

	go func(cmd, arg string) {
		command := exec.Command(cmd, arg)

		out, err := command.Output()
		if err != nil {
			ch <- err.Error()
			return
		}

		ch <- string(out)
	}(cmd, arg)

	fmt.Printf("result: %v", <-ch)
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "[error] %v\n", err)
}
