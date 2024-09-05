package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"slices"
)

func main() {
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

	entries = FoldersAtBeginning(entries)

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		fmt.Print(FileColor(info, e.Name()))
		fmt.Print(e.Name())
		fmt.Println("\x1b[0m") // Reset
	}
}
