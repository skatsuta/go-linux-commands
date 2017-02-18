package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	endMark = "$\n"
	tabMark = "^I"
)

var (
	showAll  = flag.Bool("A", false, "equivalent to -ET")
	showEnds = flag.Bool("E", false, "display $ at end of each line")
	showTabs = flag.Bool("T", false, "display TAB characters as "+string(tabMark))
)

func main() {
	flag.Parse()

	if *showAll {
		*showEnds = true
		*showTabs = true
	}

	files := flag.Args()

	if len(files) == 0 {
		// append a hyphen as stdin
		files = append(files, "-")
	}

	for _, file := range files {
		var err error
		if file == "-" {
			err = cat(os.Stdin, os.Stdout)
		} else {
			err = catFile(file)
		}

		if err != nil {
			printErr(err)
		}
	}
}

func catFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if stat, err := file.Stat(); err != nil {
		return err
	} else if stat.IsDir() {
		return fmt.Errorf("%s: Is a directory", stat.Name())
	}

	return cat(file, os.Stdout)
}

func cat(rd io.Reader, wt io.Writer) error {
	r := bufio.NewReader(rd)
	w := bufio.NewWriter(wt)

	for {
		c, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		switch {
		case c == '\n' && *showEnds:
			_, err = w.WriteString(endMark)
		case c == '\t' && *showTabs:
			_, err = w.WriteString(tabMark)
		default:
			err = w.WriteByte(c)
		}

		if err != nil {
			w.Flush()
			return err
		}

		if c == '\n' {
			w.Flush()
		}
	}

	return w.Flush()
}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
}
