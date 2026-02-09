package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"strconv"
	"strings"
)

// Trims the last decimals up to maxDecimals, does nothing if maxDecimals is less than 0, e.g -1
func trimLastDecimals(numberString string, maxDecimals int) string {
	if maxDecimals < 0 {
		return numberString
	}

	dotIndex := strings.Index(numberString, ".")
	if dotIndex == -1 {
		return numberString
	}

	return numberString[:min(len(numberString), dotIndex+maxDecimals+1)]
}

// If maxDecimals is less than 0, e.g -1, we show the exact size down to the byte
// https://en.wikipedia.org/wiki/Byte#Multiple-byte_units
func BytesToHumanReadableUnitString(bytes uint64, maxDecimals int) string {
	unitValues := []float64{
		math.Pow(10, 3),
		math.Pow(10, 6),
		math.Pow(10, 9),
		math.Pow(10, 12),
		math.Pow(10, 15),
		math.Pow(10, 18), // Largest unit that fits in 64 bits
	}

	unitStrings := []string{
		"kB",
		"MB",
		"GB",
		"TB",
		"PB",
		"EB",
	}

	if bytes < uint64(unitValues[0]) {
		return strconv.FormatUint(bytes, 10) + " B"
	}

	for i, v := range unitValues {
		if bytes >= uint64(v) {
			continue
		}

		lastIndex := max(0, i-1)
		return trimLastDecimals(strconv.FormatFloat(float64(bytes)/unitValues[lastIndex], 'f', -1, 64), maxDecimals) + " " + unitStrings[lastIndex]
	}

	return trimLastDecimals(strconv.FormatFloat(float64(bytes)/unitValues[len(unitValues)-1], 'f', -1, 64), maxDecimals) + " " + unitStrings[len(unitStrings)-1]
}

type DirStats struct {
	indexForCSV int

	totalEntriesIncludingFolders int
	numErrors                    int

	sumPathLen int
	maxPathLen int
	minPathLen int
	avgPathLen int

	numFolders               int
	numFiles                 int
	numSymlinks              int
	numExecutables           int
	numExecutablesThatAreELF int

	totalFileSize int

	numHiddenFiles int // Files starting with a period '.'
}

func GetDirStats(path string) (DirStats, []DirStats, error) {
	var out DirStats
	out.minPathLen = 2147483647 // Max 32-bit signed int (should prob be 64-bit but whatever)
	var graphOut []DirStats
	interval := 10
	i := 0

	err := myWalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if i%interval == 0 {
			out.indexForCSV = i

			outModified := out
			if outModified.minPathLen == 2147483647 { // Max 32-bit signed int
				outModified.minPathLen = 0
			}
			graphOut = append(graphOut, outModified)
		}
		i++

		if d.Name() == "." {
			return nil
		}

		if err != nil {
			out.numErrors++
			return nil
		}

		out.totalEntriesIncludingFolders++

		if d.IsDir() {
			out.numFolders++
		} else if d.Type().IsRegular() {
			out.numFiles++
			if d.Type().Perm()&0111 != 0 { // If any (owner, group, other) bits are executable
				out.numExecutables++
				// TODO: Check for ELF header for out.numExecutablesThatAreELF
			}
		} else {
			out.numSymlinks++ // FIXME: Not sure if this is correct
		}

		out.maxPathLen = max(out.maxPathLen, len(path))
		out.minPathLen = min(out.minPathLen, len(path))
		out.sumPathLen += len(path)
		out.avgPathLen = out.sumPathLen / out.totalEntriesIncludingFolders

		if !d.IsDir() {
			stat, err := os.Lstat(path)
			if err == nil {
				out.totalFileSize += int(stat.Size()) // Includes the size of symlinks (their target path length in bytes)
			}
		}

		if strings.HasPrefix(d.Name(), ".") {
			out.numHiddenFiles++
		}

		return nil
	})
	if err != nil {
		return out, graphOut, err
	}

	return out, graphOut, nil
}

func colorOfNumber(n int) string {
	if n <= 0 {
		return "\x1b[0m"
	} else {
		return "\x1b[1;32m" // Green
	}
}

func colorNumber(n int) string {
	return colorOfNumber(n) + strconv.Itoa(n) + "\x1b[0m"
}

func main() {
	stats, graph, err := GetDirStats(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reset := "\x1b[0m"
	gray := "\x1b[0;37m"

	fmt.Println(gray+"total entries including folders:", colorNumber(stats.totalEntriesIncludingFolders))
	fmt.Println(gray+"number of errors:", colorNumber(stats.numErrors))
	fmt.Println("")

	fmt.Println(gray+"sum path length:", colorNumber(stats.sumPathLen))
	fmt.Println(gray+"max path length:", colorNumber(stats.maxPathLen))
	fmt.Println(gray+"min path length:", colorNumber(stats.minPathLen))
	fmt.Println(gray+"avg path length:", colorNumber(stats.avgPathLen))
	fmt.Println("")

	fmt.Println(gray+"folders:     ", colorNumber(stats.numFolders))
	fmt.Println(gray+"files:       ", colorNumber(stats.numFiles))
	fmt.Println(gray+"symlinks:    ", colorNumber(stats.numSymlinks))
	fmt.Println(gray+"hidden files:", colorNumber(stats.numHiddenFiles))
	fmt.Println("")

	fmt.Println(gray+"total file size:", colorOfNumber(stats.totalFileSize), BytesToHumanReadableUnitString(uint64(stats.totalFileSize), -1), reset)

	// TODO:
	// executables:
	// # of which are ELF executables: (header)

	if false {
		// CSV
		fmt.Println("index,# entries incl. folders,# errors,max path len,min path len,avg path len,sum path len,# folders,# files,# symlinks,# hidden files,total file size")
		for _, e := range graph {
			fmt.Print(e.indexForCSV, ",")
			fmt.Print(e.totalEntriesIncludingFolders, ",")
			fmt.Print(e.numErrors, ",")
			fmt.Print(e.maxPathLen, ",")
			fmt.Print(e.minPathLen, ",")
			fmt.Print(e.avgPathLen, ",")
			fmt.Print(e.sumPathLen, ",")
			fmt.Print(e.numFolders, ",")
			fmt.Print(e.numFiles, ",")
			fmt.Print(e.numSymlinks, ",")
			fmt.Print(e.numHiddenFiles, ",")

			fmt.Println(e.totalFileSize)
		}
	}
}
