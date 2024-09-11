package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kivattt/getopt"
)

func main() {
	sevenBit := flag.Bool("seven-bit", false, "Only print up to byte value 127")
	printable := flag.Bool("printable", false, "Only printable ASCII values from 0x20 to 0x7e")
	help := flag.Bool("help", false, "Print help")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("ascii", flag.ExitOnError)
	getopt.Aliases(
		"s", "seven-bit",
		"p", "printable",
		"h", "help",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *help {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS]")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	minimum := 0
	maximum := 256
	if *sevenBit {
		maximum = 128
	}

	if *printable {
		minimum = 0x20
		maximum = 127
	}

	byteSlice := []byte{}
	for i := minimum; i < maximum; i++ {
		byteSlice = append(byteSlice, byte(i))
	}
	os.Stdout.Write(byteSlice)
}
