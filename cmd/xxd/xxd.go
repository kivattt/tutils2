package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kivattt/getopt"
	"golang.org/x/term"
)

func isPrintableASCIIRange(c byte) bool {
	return c >= 0x20 && c <= 0x7e
}

func charColor(c byte) string {
	if c == 0 {
		return "\x1b[1;37m" // White
	}

	if c == ' ' || c == 0xff {
		return "\x1b[1;34m" // Blue
	}

	if strings.ContainsRune("\n\r\t", rune(c)) {
		return "\x1b[1;33m" // Yellow
	}

	if isPrintableASCIIRange(c) {
		return "\x1b[1;32m" // Green
	}

	return "\x1b[1;31m" // Red
}

func coloredText(bytes []byte, colorsEnabled bool) string {
	var builder strings.Builder
	for _, b := range bytes {
		if colorsEnabled {
			builder.WriteString(charColor(b))
		}

		if isPrintableASCIIRange(b) {
			builder.WriteByte(b)
		} else {
			builder.WriteByte('.')
		}
	}

	if colorsEnabled {
		builder.WriteString("\x1b[0m") // Reset
	}
	return builder.String()
}

func leadingZeroesGray(str string) string {
	firstZero := strings.IndexFunc(str, func(r rune) bool {
		return r != '0'
	})
	firstZero = min(7, firstZero)
	return "\x1b[0;37m" + str[:firstZero] + "\x1b[0m" + str[firstZero:]
}

// Returns number of lines output
func handleBuf(buf []byte, size, width, nthLineOutput int, colorsEnabled, decimalInsteadOfHex bool) int {
	nLinesOutput := 0

	for i := 0; i <= size; i += width {
		if i >= size {
			break
		} else {
			formatString := "%08x: "
			if decimalInsteadOfHex {
				formatString = "%08d: "
			}

			if colorsEnabled {
				fmt.Print(leadingZeroesGray(fmt.Sprintf(formatString, i + nthLineOutput*width)))
			} else {
				fmt.Print(fmt.Sprintf(formatString, i + nthLineOutput*width))
			}
		}

		nCharsPrinted := 0
		for j := i; j < min(size, i+width); j++ {
			if colorsEnabled {
				fmt.Printf("%s%02x", charColor(buf[j]), buf[j])
			} else {
				fmt.Printf("%02x", buf[j])
			}
			nCharsPrinted += 2

			if j % 2 == 1 && j != min(size, i+width) - 1 {
				fmt.Print(" ")
				nCharsPrinted++
			}
		}

		if colorsEnabled {
			fmt.Print("\x1b[0m")
		}
		fmt.Print(strings.Repeat(" ", max(0, (width/2 + width*2)-nCharsPrinted)) + " " + coloredText(buf[i:min(size, i+16)], colorsEnabled))
		fmt.Println()
		nLinesOutput++
	}

	return nLinesOutput
}

func main() {
	help := flag.Bool("help", false, "display this help and exit")
	decimal := flag.Bool("decimal", false, "show offset in decimal instead of hex")
	color := flag.String("color", "auto", "colorize the output [auto, always, never]")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init("xxd", flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"d", "decimal",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(0)
	}

	if *help {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS] [FILES]")
		fmt.Println("Show as hex dump")
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

	buf := make([]byte, 32) // This has to be like, a power of two more than or equal to 32
	width := 16 // If you change this it doesn't output correctly...

	nthLineOutput := 0

	// Read files
	if len(getopt.CommandLine.Args()) > 0 {
		for _, path := range getopt.CommandLine.Args() {
			nthLineOutput = 0
			f, err := os.Open(path)
			if err != nil {
				continue
			}

			for {
				n, err := io.ReadFull(f, buf)
				nthLineOutput += handleBuf(buf, n, width, nthLineOutput, colorToUse != "never", *decimal)

				// End of file
				if err != nil {
					break
				}
			}

			f.Close()
		}

		if len(getopt.CommandLine.Args()) > 1 {
			if colorToUse != "never" {
				os.Stderr.WriteString("\x1b[1;31m") // Red
			}

			os.Stderr.WriteString("Using multiple files with xxd does not internally concatenate them,\n")
			os.Stderr.WriteString("and produces non-standard output\n")

			if colorToUse != "never" {
				os.Stderr.WriteString("\x1b[0m") // Reset
			}
		}

		os.Exit(0)
	}

	stat, _ := os.Stdin.Stat()
	// Not piped input
	if stat.Mode() & os.ModeCharDevice != 0 {
		for {
			n, err := os.Stdin.Read(buf)
			nthLineOutput += handleBuf(buf, n, width, nthLineOutput, colorToUse != "never", *decimal)

			// End of file
			if err != nil {
				break
			}
		}
		os.Exit(0)
	}

	// Piped input
	for {
		n, err := io.ReadFull(os.Stdin, buf)
		nthLineOutput += handleBuf(buf, n, width, nthLineOutput, colorToUse != "never", *decimal)

		// End of file
		if err != nil {
			break
		}
	}
}
