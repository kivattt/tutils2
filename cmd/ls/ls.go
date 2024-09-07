package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kivattt/getopt"
)

var validSortByValues = [...]string{
	"none",
	"modified",
}

func main() {
	h := flag.Bool("help", false, "display this help and exit")
	all := flag.Bool("all", false, "Show files starting with '.'")
	sortBy := flag.String("sort-by", "none", "sort files ("+strings.Join(validSortByValues[:], ", ")+")")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("ls", flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"a", "all",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *h {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS]")
		fmt.Println("List files")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	path, err := os.Getwd()
	if err != nil {
		path = os.Getenv("PWD")
		if path == "" {
			log.Fatal("PWD environment variable empty, unable to determine current working directory")
		}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal("Failed to read directory '" + path + "'")
	}

	switch *sortBy {
	case "modified":
		slices.SortStableFunc(entries, func(a, b fs.DirEntry) int {
			aInfo, aErr := a.Info()
			bInfo, bErr := b.Info()
			if aErr != nil || bErr != nil {
				return 0
			}

			if aInfo.ModTime().Before(bInfo.ModTime()) {
				return -1
			}

			if aInfo.ModTime().Equal(bInfo.ModTime()) {
				return 0
			}

			return 1
		})
	case "none":
	default:
		fmt.Fprintln(os.Stderr, "Invalid sortBy value \""+*sortBy+"\"")
		fmt.Fprintln(os.Stderr, "Valid values: "+strings.Join(validSortByValues[:], ", "))
		os.Exit(1)
	}

	entries = FoldersAtBeginning(entries)

	for _, e := range entries {
		if !*all && strings.HasPrefix(e.Name(), ".") {
			continue
		}

		info, err := e.Info()
		if err != nil {
			continue
		}
		fmt.Print(FileColor(info, e.Name()))
		fmt.Print(e.Name())
		fmt.Println("\x1b[0m") // Reset
	}
}
