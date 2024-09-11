package main

import (
	"fmt"
	"os"
)

func main() {
	buf := make([]byte, 1024)
	for {
		n, err := os.Stdin.Read(buf)

		for i, c := range buf[:n] {
			if i%16 == 0 {
				if i != 0 {
					fmt.Print("\n")
				}
				fmt.Printf("%08x: ", i)
			} else if i%2 == 0 {
				fmt.Print(" ")
			}

			fmt.Printf("%02x", c)
		}

		// End of file
		if err != nil {
			break
		}
	}

	fmt.Println()
}
