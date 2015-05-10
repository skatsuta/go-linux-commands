package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {

	for i := 1; i < len(os.Args); i++ {
		fp, err := os.Open(os.Args[i])
		if err != nil {
			die(err.Error())
		}
		defer func() {
			if e := fp.Close(); e != nil {
				die(e.Error())
			}
		}()

		r := bufio.NewReader(fp)
		cnt := 0

		for {
			c, err := r.ReadByte()
			if c == '\n' {
				cnt++
			}
			if err == io.EOF {
				break
			}
		}

		fmt.Printf("%s: %d lines\n", os.Args[i], cnt)
	}

}

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
