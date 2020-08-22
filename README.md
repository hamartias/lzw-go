### LZW Compression/Decompression

[Lempel-Ziv-Welch](https://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Welch) is a straight-forward dictionary based compression algorithm.

##### Usage
For a usage example of the lzw package see
[cmd/main.go](https://github.com/davidcrosby/lzw-go/blob/master/cmd/main.go)

Showing compression/decompression and compressed size for [a big file.](https://norvig.com/big.txt)

    $ go run lzw.go -c test_data/big.txt big.out
    $ go run lzw.go -d big.out big.decompressed
    $ du test_data/big.txt
    6340	test_data/big.txt
    $ du big.out
    3300	big.out
    $ diff big.decompressed test_data/big.txt
    $ 

