package main

import (
	"fmt"
	"os"

	"github.com/davidcrosby/lzw-go"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: `go run lzw.go -(d|c) (input_file) (output_file)`")
		return
	}
	opt := os.Args[1]
	inputFname := os.Args[2]
	outputFname := os.Args[3]
	if opt == "-c" {
		lzw.Compress(inputFname, outputFname)
	} else if opt == "-d" {
		lzw.Decompress(inputFname, outputFname)
	}
}
