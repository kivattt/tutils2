package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/kivattt/getopt"
	"golang.org/x/term"
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

func printError(msg string, colorEnabled bool) {
	if colorEnabled {
		os.Stderr.WriteString("\x1b[1;31m") // Red
	}
	os.Stderr.WriteString(msg + "\n")
	if colorEnabled {
		os.Stderr.WriteString("\x1b[0m") // Reset
	}
}

func main() {
	help := flag.Bool("help", false, "display this help and exit")
	all := flag.Bool("all", false, "show hidden files starting with '.'")
	directoriesFirst := flag.Bool("directories-first", false, "show directories first")
	directory := flag.Bool("directory", false, "list directories themselves, not their contents")
	sortBy := flag.String("sort-by", "none", "sort files ("+strings.Join(validSortByValues[:], ", ")+")")
	summary := flag.Bool("summary", false, "folder stats")
	color := flag.String("color", "auto", "colorize the output [auto, always, never]")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("ls", flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"a", "all",
		"d", "directory",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *help {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS] [FILES]")
		fmt.Println("List files")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	colorToUse := *color
	if colorToUse == "auto" {
		if !term.IsTerminal(int(os.Stdout.Fd())) {
			colorToUse = "never" // Output is piped, don't colorize the output
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = os.Getenv("PWD")
		if cwd == "" {
			printError("PWD environment variable empty, unable to determine current working directory", colorToUse != "never")
			os.Exit(1)
		}
	}

	var paths []string
	if len(getopt.CommandLine.Args()) == 0 {
		paths = []string{cwd}
	} else {
		paths = getopt.CommandLine.Args()
	}

	var allEntries []fs.DirEntry

	for _, path := range paths {
		stat, err := os.Lstat(path)
		if err != nil {
			printError("Failed to stat: '"+path+"'", colorToUse != "never")
			continue
		}

		if !stat.IsDir() || *directory {
			allEntries = append(allEntries, fs.FileInfoToDirEntry(stat))
			continue
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			printError("Failed to read directory '"+path+"'", colorToUse != "never")
			continue
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
		os.Exit(0)
	}

	SortEntries(&allEntries, *sortBy)

	if *directoriesFirst {
		allEntries = FoldersAtBeginning(allEntries)
	}

	for _, e := range allEntries {
		if !*all && strings.HasPrefix(e.Name(), ".") {
			continue
		}

		info, err := e.Info()
		if err != nil {
			printError("Failed to stat: '"+e.Name()+"'", colorToUse != "never")
			continue
		}

		if colorToUse != "never" {
			fmt.Print(FileColor(info, e.Name()))
		}

		fmt.Println(e.Name())

		if colorToUse != "never" {
			fmt.Print("\x1b[0m")
		}
	}
}
