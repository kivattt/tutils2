package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// This function taken from util.go in https://github.com/kivattt/fen
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

// This function taken from util.go in https://github.com/kivattt/fen
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

func PathWithEndSeparator(path string) string {
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return path
	}

	return path + string(os.PathSeparator)
}

func usage(programName string) {
	fmt.Println("Usage: " + programName + " [bytes number]")
	fmt.Println("Show bytes in human-readable size")
}

func main() {
	if len(os.Args) < 2 {
		usage(os.Args[0])
		os.Exit(0)
	}

	nBytes, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage(os.Args[0])
		os.Exit(1)
	}

	fmt.Println(nBytes, "bytes is equivalent to", BytesToHumanReadableUnitString(uint64(nBytes), 3))
}
