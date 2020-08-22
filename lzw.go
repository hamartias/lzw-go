package lzw

import (
	"io"
	"os"

	"github.com/dgryski/go-bitstream"
)

var codeWidth = 12
var maxCode int64 = 4095

func check(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func Compress(inputFname, outputFname string) {
	table := make(map[string]int64)
	// initialize with extended ascii
	for char := 0; char < 256; char++ {
		table[string(rune(char))] = int64(char)
	}

	input, err := os.Open(inputFname)
	check(err)

	output, err := os.Create(outputFname)
	check(err)

	reader := bitstream.NewReader(input)
	writer := bitstream.NewWriter(output)

	var nextCodeToAdd int64 = 256
	var runningString []byte
	toAdd, err := reader.ReadByte()
	check(err)
	runningString = append(runningString, toAdd)

	for {
		char, err := reader.ReadByte()
		if err == io.EOF {
			finalCode := table[string(runningString)]

			// check for empty file
			if finalCode != 0 {
				errw := writer.WriteBits(uint64(finalCode), codeWidth)
				check(errw)
				errw = writer.Flush(false)
				check(errw)
			}
			break
		}
		check(err)
		withChar := string(runningString) + string(char)
		if _, found := table[withChar]; !found {
			code := table[string(runningString)]
			err := writer.WriteBits(uint64(code), codeWidth)
			check(err)

			if nextCodeToAdd < maxCode {
				table[withChar] = nextCodeToAdd
				nextCodeToAdd++
			}
			runningString = nil
		}
		runningString = append(runningString, char)
	}

}

func Decompress(inputFname string, outputFname string) {
	table := make(map[int64]string)
	// initialize to extended ascii
	for code := 0; code < 256; code++ {
		table[int64(code)] = string(rune(code))
	}

	input, err := os.Open(inputFname)
	check(err)

	output, err := os.Create(outputFname)
	check(err)

	reader := bitstream.NewReader(input)
	writer := bitstream.NewWriter(output)

	var nextCodeToAdd int64 = 256

	lastCodeUns, err := reader.ReadBits(codeWidth)
	check(err)

	var lastCode int64 = int64(lastCodeUns)

	// empty file
	if lastCode == 0 {
		return
	}

	werr := writer.WriteByte(byte(table[lastCode][0]))
	check(werr)
	oldCode := lastCode

	for {
		codeUnsigned, err := reader.ReadBits(codeWidth)
		if err == io.EOF {
			break
		}
		check(err)
		var code int64 = int64(codeUnsigned)

		var outputString string
		if _, found := table[code]; !found {
			lastString := table[oldCode]
			outputString = lastString + string(lastString[0])
		} else {
			outputString = table[code]
		}
		for _, val := range []byte(outputString) {
			werr := writer.WriteByte(val)
			check(werr)
		}
		if nextCodeToAdd < maxCode {
			nextStringToAdd := table[oldCode] + string(outputString[0])
			table[nextCodeToAdd] = nextStringToAdd
			nextCodeToAdd++
		}
		oldCode = code
	}
}
