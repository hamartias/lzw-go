package lzw

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/dgryski/go-bitstream"
)

func check(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

// Config Input/Output streams and compression parameters
type Config struct {
	r          *bitstream.BitReader
	w          *bitstream.BitWriter
	codeWidth  int
	maxCode    int64
	dictionary map[string]int64
}

// SetupConfig Configuration struct for compression/decompression
func SetupConfig(
	input io.Reader,
	output io.Writer,
	codeWidth int,
	dictionary map[string]int64,
) (*Config, error) {
	conf := new(Config)
	conf.r = bitstream.NewReader(input)
	conf.w = bitstream.NewWriter(output)
	conf.codeWidth = codeWidth
	conf.maxCode = int64(math.Exp2(float64(codeWidth)))
	conf.dictionary = dictionary
	err := error(nil)
	for key, code := range dictionary {
		if code > conf.maxCode {
			msg := fmt.Sprintf("Value for %s is too large for this code width.", key)
			err = errors.New(msg)
		}
	}

	return conf, err
}

// Compress all input from input and write it to output
func Compress(conf *Config) {
	table := conf.dictionary

	var nextCodeToAdd int64 = 256
	var runningString []byte
	toAdd, err := conf.r.ReadByte()
	check(err)
	runningString = append(runningString, toAdd)

	for {
		char, err := conf.r.ReadByte()
		if err == io.EOF {
			finalCode := table[string(runningString)]

			// check for empty file
			if finalCode != 0 {
				errw := conf.w.WriteBits(uint64(finalCode), conf.codeWidth)
				check(errw)
				errw = conf.w.Flush(false)
				check(errw)
			}
			break
		}
		check(err)
		withChar := string(runningString) + string(char)
		if _, found := table[withChar]; !found {
			code := table[string(runningString)]
			err := conf.w.WriteBits(uint64(code), conf.codeWidth)
			check(err)

			if nextCodeToAdd < conf.maxCode {
				table[withChar] = nextCodeToAdd
				nextCodeToAdd++
			}
			runningString = nil
		}
		runningString = append(runningString, char)
	}

}

// Decompress input from input into output
func Decompress(conf *Config) {
	// Inverse config dictionary
	table := make(map[int64]string)
	for key, val := range conf.dictionary {
		table[val] = key
	}

	var nextCodeToAdd int64 = 256

	lastCodeUns, err := conf.r.ReadBits(conf.codeWidth)
	check(err)
	var lastCode int64 = int64(lastCodeUns)

	// empty file
	if lastCode == 0 {
		return
	}

	werr := conf.w.WriteByte(byte(table[lastCode][0]))
	check(werr)
	oldCode := lastCode

	for {
		codeUnsigned, err := conf.r.ReadBits(conf.codeWidth)
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
			werr := conf.w.WriteByte(val)
			check(werr)
		}
		if nextCodeToAdd < conf.maxCode {
			nextStringToAdd := table[oldCode] + string(outputString[0])
			table[nextCodeToAdd] = nextStringToAdd
			nextCodeToAdd++
		}
		oldCode = code
	}
}
