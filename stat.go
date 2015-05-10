package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type stat struct {
	fi os.FileInfo
}

func (s *stat) String() string {
	format := `
type	%v
size	%v
mode	%v
modTime	%v
sys	%#v
`

	return fmt.Sprintf(format, s.fi.Name(), s.fi.Size(), s.fi.Mode(), s.fi.ModTime(), s.fi.Sys())
}

func newStat(fi os.FileInfo) *stat {
	return &stat{fi: fi}
}

func main() {
	if len(os.Args) != 2 {
		perror(fmt.Errorf("%s: wrong arguments\n", os.Args[0]))
		return
	}

	fi, err := os.Lstat(os.Args[1])
	if err != nil {
		perror(err)
		return
	}

	fmt.Println(newStat(fi))
}

func doLs(path string) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		perror(err)
		return
	}

	for _, fi := range fis {
		fmt.Printf("%s\n", fi.Name())
	}
}

func perror(err error) {
	fmt.Fprintf(os.Stderr, "[error] %v\n", err)
}
