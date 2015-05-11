package main

import (
	"fmt"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println(u.Username)
}
