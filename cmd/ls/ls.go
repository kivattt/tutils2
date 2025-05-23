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
	"github.com/kivattt/gogitstatus"
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
	demo := flag.Bool("demo", false, "show all the file colors")
	gitStatus := flag.Bool("git-status", false, "highlight changed/untracked files from current working directory")
	gitStatusDetailed := flag.Bool("git-status-detailed", false, "show more info about changed/untracked files")

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

	if *demo {
		for key, val := range colors {
			if colorToUse != "never" {
				fmt.Print(val)
			}
			fmt.Println(key)
			if colorToUse != "never" {
				fmt.Print("\x1b[0m") // Reset color
			}
		}

		os.Exit(0)
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
	longestEntryBasename := 0

	for _, path := range paths {
		stat, err := os.Lstat(path)
		if err != nil {
			printError("Failed to stat: '"+path+"'", colorToUse != "never")
			continue
		}

		if !stat.IsDir() || *directory {
			longestEntryBasename = max(longestEntryBasename, len(stat.Name()))
			allEntries = append(allEntries, fs.FileInfoToDirEntry(stat))
			continue
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			printError("Failed to read directory '"+path+"'", colorToUse != "never")
			continue
		}

		for _, entry := range entries {
			longestEntryBasename = max(longestEntryBasename, len(entry.Name()))
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

	var changedOrUntracked map[string]gogitstatus.ChangedFile
	if *gitStatus {
		changedOrUntracked, _ = gogitstatus.Status(".")
		changedOrUntracked = gogitstatus.IncludingDirectories(changedOrUntracked)
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

		// FIXME: Check correct parent folder first, e.g. with filepath.Rel() ?
		changedFile, ok := changedOrUntracked[e.Name()]

		if colorToUse != "never" {
			if ok {
				fmt.Print("\x1b[0;31m") // Red background
			} else {
				fmt.Print(FileColor(info, e.Name()))
			}
		}

		if ok && *gitStatusDetailed {
			fmt.Print(e.Name())
			if colorToUse != "never" {
				fmt.Print("\x1b[0m")
			}

			nSpaces := 1 + longestEntryBasename - len(e.Name())
			spaces := strings.Repeat(" ", nSpaces)

			whatChanged := gogitstatus.WhatChangedToString(changedFile.WhatChanged)
			nSpaces2nd := 1 + len("OWNER_CHANGED") - len(whatChanged) // "OWNER_CHANGED" is the longest possible
			spaces2nd := strings.Repeat(" ", nSpaces2nd)

			if changedFile.Untracked {
				fmt.Print(spaces + whatChanged + spaces2nd)
				if colorToUse != "never" {
					fmt.Print("\x1b[0;31m") // Red
				}
				fmt.Print("untracked")
				if colorToUse != "never" {
					fmt.Print("\x1b[0m")
				}
				fmt.Println()
			} else {
				fmt.Print(spaces + whatChanged + spaces2nd)
				if colorToUse != "never" {
					fmt.Print("\x1b[38;2;254;229;65m") // Yellow
				}
				fmt.Print("unstaged")
				if colorToUse != "never" {
					fmt.Print("\x1b[0m")
				}
				fmt.Println()
			}
		} else {
			fmt.Print(e.Name())
			if colorToUse != "never" {
				fmt.Print("\x1b[0m")
			}
			fmt.Println()
		}

		if colorToUse != "never" {
			fmt.Print("\x1b[0m")
		}
	}
}
