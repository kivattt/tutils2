package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/kivattt/getopt"
	"golang.org/x/term"
)

func printPathError(path string, colorsEnabled bool) {
	if colorsEnabled {
		os.Stderr.WriteString("\x1b[1;31m") // Red
	}
	os.Stderr.WriteString("No such file: " + path + "\n")
	if colorsEnabled {
		os.Stderr.WriteString("\x1b[0m")
	}
}

func tabsToSpaces(buf []byte) []byte {
	ret := buf
	for i, c := range ret {
		if c == '\t' {
			ret[i] = ' '
			ret = slices.Insert(ret, i, []byte{' ', ' ', ' '}...)
		}
	}
	return ret
}

func main() {
	help := flag.Bool("help", false, "display this help and exit")
	color := flag.String("color", "auto", "colorize stderr messages [auto, always, never]")
	fourSpaces := flag.Bool("four", false, "turn tabs into 4 spaces")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("cat", flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"4", "four",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *help {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS] [FILES]")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	colorToUse := *color
	if colorToUse == "auto" {
		if !term.IsTerminal(int(os.Stderr.Fd())) {
			colorToUse = "never" // Output is piped, don't colorize our error messages
		}
	}

	if len(getopt.CommandLine.Args()) > 0 {
		for _, path := range getopt.CommandLine.Args() {
			f, err := os.Open(path)
			if err != nil {
				printPathError(path, colorToUse != "never")
				continue
			}

			buf := make([]byte, 512)
			for {
				n, err := f.Read(buf)

				if *fourSpaces {
					os.Stdout.Write(tabsToSpaces(buf[:n]))
				} else {
					os.Stdout.Write(buf[:n])
				}

				// End of file
				if err != nil {
					break
				}
			}

			f.Close()
		}

		os.Exit(0)
	}

	// Stdin input
	buf := make([]byte, 512)
	for {
		n, err := os.Stdin.Read(buf)

		if *fourSpaces {
			os.Stdout.Write(tabsToSpaces(buf[:n]))
		} else {
			os.Stdout.Write(buf[:n])
		}

		// End of file
		if err != nil {
			break
		}
	}
}
