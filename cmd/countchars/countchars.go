package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/kivattt/getopt"
)

func printableRangeASCII(c byte) bool {
	return c >= 0x20 && c <= 127
}

func printableChar(c byte, decimal bool) string {
	if printableRangeASCII(c) {
		return string(c)
	}

	if decimal {
		return strconv.Itoa(int(c))
	}

	return fmt.Sprintf("0x%02x", c)
}

func main() {
	decimal := flag.Bool("decimal", false, "show ASCII codes in decimal instead of hexadecimal")
	h := flag.Bool("help", false, "display this help and exit")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("countchars", flag.ExitOnError)
	getopt.Aliases(
		"d", "decimal",
		"h", "help",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *h {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS]")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	buf := make([]byte, 512)

	byteMap := make(map[byte]uint64)
	for {
		n, err := os.Stdin.Read(buf)

		for _, c := range buf[:n] {
			byteMap[c]++
		}

		if err != nil {
			break
		}
	}

	type kv struct {
		Key byte
		Value uint64
	}
	var ss []kv
	for k, v := range byteMap {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		fmt.Println(printableChar(kv.Key, *decimal) + " : " + strconv.FormatUint(kv.Value, 10))
	}
}
