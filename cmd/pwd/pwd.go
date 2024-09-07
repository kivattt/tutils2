package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		path = os.Getenv("PWD")
		if path == "" {
			log.Fatal("PWD environment variable empty, unable to determine current working directory")
		}
	}

	fmt.Println(path)
}
