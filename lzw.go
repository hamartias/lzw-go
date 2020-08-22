package lzw

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/dgryski/go-bitstream"
)

// Config encapsulates bitstreams and compression parameters
type Config struct {
	r          *bitstream.BitReader
	w          *bitstream.BitWriter
	codeWidth  int
	maxCode    int64
	dictionary map[string]int64
}

// SetupConfig sets up Config struct for compression/decompression with input
// validation
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

//Compress performs compression according to parameters in Config
func Compress(conf *Config) error {
	table := conf.dictionary

	var nextCodeToAdd int64 = 256
	var runningString []byte
	toAdd, err := conf.r.ReadByte()
	if err != nil {
		return err
	}
	runningString = append(runningString, toAdd)

	for {
		char, err := conf.r.ReadByte()
		if err == io.EOF {
			finalCode := table[string(runningString)]
			// check for empty file
			if finalCode != 0 {
				errw := conf.w.WriteBits(uint64(finalCode), conf.codeWidth)
				if errw != nil {
					return errw
				}
				errw = conf.w.Flush(false)
				if errw != nil {
					return err
				}
			}
			break
		}
		if err != nil {
			return err
		}
		withChar := string(runningString) + string(char)
		if _, found := table[withChar]; !found {
			code := table[string(runningString)]
			err := conf.w.WriteBits(uint64(code), conf.codeWidth)
			if err != nil {
				return err
			}

			if nextCodeToAdd < conf.maxCode {
				table[withChar] = nextCodeToAdd
				nextCodeToAdd++
			}
			runningString = nil
		}
		runningString = append(runningString, char)
	}
	return error(nil)
}

// Decompress performs decompression according to parameters in Config
func Decompress(conf *Config) error {
	// Inverse config dictionary
	table := make(map[int64]string)
	for key, val := range conf.dictionary {
		table[val] = key
	}

	var nextCodeToAdd int64 = 256

	lastCodeUns, err := conf.r.ReadBits(conf.codeWidth)
	if err != nil {
		return err
	}
	var lastCode int64 = int64(lastCodeUns)

	// empty file
	if lastCode == 0 {
		return error(nil)
	}

	werr := conf.w.WriteByte(byte(table[lastCode][0]))
	if werr != nil {
		return werr
	}
	oldCode := lastCode

	for {
		codeUnsigned, err := conf.r.ReadBits(conf.codeWidth)
		if err != nil {
			return err
		}
		if err == io.EOF {
			break
		}
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
			if werr != nil {
				return werr
			}
		}
		if nextCodeToAdd < conf.maxCode {
			nextStringToAdd := table[oldCode] + string(outputString[0])
			table[nextCodeToAdd] = nextStringToAdd
			nextCodeToAdd++
		}
		oldCode = code
	}
	return error(nil)
}
