package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/kivattt/getopt"
)

var validSortByValues = [...]string{
	"none",
	"modified",
}

func SortEntries(entries *[]fs.DirEntry, sortBy string) {
	switch sortBy {
	case "modified":
		slices.SortStableFunc(*entries, func(a, b fs.DirEntry) int {
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
		fmt.Fprintln(os.Stderr, "Invalid sortBy value \""+sortBy+"\"")
		fmt.Fprintln(os.Stderr, "Valid values: "+strings.Join(validSortByValues[:], ", "))
		os.Exit(1)
	}
}

func main() {
	h := flag.Bool("help", false, "display this help and exit")
	all := flag.Bool("all", false, "Show files starting with '.'")
	foldersFirst := flag.Bool("folders-first", false, "show folders first")
	summary := flag.Bool("summary", false, "Folder stats")
	sortBy := flag.String("sort-by", "none", "sort files ("+strings.Join(validSortByValues[:], ", ")+")")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("ls", flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"a", "all",
		"f", "folders-first",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *h {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS] [FILES]")
		fmt.Println("List files")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = os.Getenv("PWD")
		if cwd == "" {
			log.Fatal("PWD environment variable empty, unable to determine current working directory")
		}
	}

	var allEntries []fs.DirEntry

	var paths []string
	if len(getopt.CommandLine.Args()) == 0 {
		paths = []string{cwd}
	} else {
		for _, e := range getopt.CommandLine.Args() {
			stat, err := os.Stat(e)
			if err != nil || !stat.IsDir() {
				eAbs, err := filepath.Abs(e)
				if err == nil && !slices.ContainsFunc(allEntries, func(e fs.DirEntry) bool {
					// Only checks base, bad
					return filepath.Base(eAbs) == e.Name()
				}) {
					allEntries = append(allEntries, fs.FileInfoToDirEntry(stat))
				}
				continue
			}

			isDuplicate := slices.ContainsFunc(paths, func(path string) bool {
				pathAbs, pathErr := filepath.Abs(path)
				eAbs, eErr := filepath.Abs(e)
				if pathErr != nil || eErr != nil {
					return false
				}

				return pathAbs == eAbs
			})

			if !isDuplicate {
				paths = append(paths, e)
			}
		}
	}

	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			log.Fatal("Failed to read directory '" + path + "'")
		}

		allEntries = append(allEntries, entries...)
	}

	// Just folder stats, don't need any sorting
	if *summary {
		folderCount := 0
		fileCount := 0
		hiddenFolderCount := 0
		hiddenFileCount := 0

		for _, e := range allEntries {
			isHidden := strings.HasPrefix(e.Name(), ".")

			if e.IsDir() {
				if isHidden {
					hiddenFolderCount++
				} else {
					folderCount++
				}
			} else {
				if isHidden {
					hiddenFileCount++
				} else {
					fileCount++
				}
			}
		}

		fmt.Print(folderCount, " folders")
		if hiddenFolderCount > 0 {
			fmt.Print(" (" + strconv.Itoa(hiddenFolderCount) + " hidden)")
		}
		fmt.Println()

		fmt.Print(fileCount, " files")
		if hiddenFileCount > 0 {
			fmt.Print("   (" + strconv.Itoa(hiddenFileCount) + " hidden)")
		}
		fmt.Println()

		fmt.Println(folderCount+hiddenFolderCount+fileCount+hiddenFileCount, "total")
	}

	SortEntries(&allEntries, *sortBy)

	if *foldersFirst {
		allEntries = FoldersAtBeginning(allEntries)
	}

	for _, e := range allEntries {
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
