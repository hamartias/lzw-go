package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davidcrosby/lzw-go"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage:\n `go run (path2command) -(d|c) (input_file) (output_file)`\n")
		return
	}
	opt := os.Args[1]
	inputFname := os.Args[2]
	outputFname := os.Args[3]

	input, err := os.Open(inputFname)
	if err != nil {
		panic(err)
	}
	output, err := os.Create(outputFname)
	if err != nil {
		panic(err)
	}

	table := make(map[string]int64)
	for code := 0; code < 256; code++ {
		table[string(rune(code))] = int64(code)
	}

	conf, _ := lzw.SetupConfig(input, output, 12, table)
	if opt == "-c" {
		err := lzw.Compress(conf)
		if err != nil {
			log.Fatal(err)
		}
	} else if opt == "-d" {
		err := lzw.Decompress(conf)
		if err != nil {
			log.Fatal(err)
		}
	}
}
