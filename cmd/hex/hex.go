package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kivattt/getopt"
)

const hexLookup = "0123456789abcdef"

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
	decode := flag.Bool("decode", false, "decode hexadecimal")
	help := flag.Bool("help", false, "display this help and exit")
	nonewline := flag.Bool("nonewline", false, "don't output trailing newline")
	noignore := flag.Bool("noignore", false, "if invalid 2-byte hex code found during decode, exit with code 1")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("hex", flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"n", "nonewline",
		"d", "decode",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *help {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS] [FILES]")
		fmt.Println("Encode/decode hexadecimal")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	invalidHexCodeFound := false
	anyData := false
	buf := make([]byte, 512)
	if *decode {
		bit := false
		var lastByte byte
		for {
			n, err := os.Stdin.Read(buf)

			for _, currentByte := range buf[:n] {
				if bit {
					left := strings.IndexByte(hexLookup, lastByte)
					right := strings.IndexByte(hexLookup, currentByte)
					if left != -1 && right != -1 {
						anyData = true
						os.Stdout.Write([]byte{byte(left<<4 | right)})
					} else {
						invalidHexCodeFound = true
						if *noignore {
							break
						}
					}
				}
				bit = !bit
				lastByte = currentByte
			}

			// End of file
			if err != nil {
				break
			}
		}
	} else {
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				anyData = true
			}

			for _, c := range buf[:n] {
				os.Stdout.Write([]byte{hexLookup[c>>4], hexLookup[c&0xf]})
			}

			// End of file
			if err != nil {
				break
			}
		}
	}

	if !*nonewline && anyData {
		fmt.Println()
	}

	if invalidHexCodeFound {
		printError("Invalid hex code found in string", true) // FIXME: Color
		if *noignore {
			os.Exit(1)
		}
	}
}
