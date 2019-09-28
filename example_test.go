package compints_test

import (
	"fmt"
	"github.com/shibukawa/compints"
	"reflect"
)

var input = []uint32{1, 2, 3, 4, 5, 6, 6, 8, 10}

func ExampleCompress() {
	compressed, cbytes, length := compints.Compress(input, true)
	output, err := compints.Decompress(compressed, cbytes, length, true)
	fmt.Printf("equal: %v\n", reflect.DeepEqual(input, output))
	fmt.Printf("error: %v\n", err)
	// Output:
	// equal: true
	// error: <nil>
}

func ExampleCompressToBytes() {
	compressed := compints.CompressToBytes(input, false)
	output, err := compints.DecompressFromBytes(compressed, false)
	fmt.Printf("equal: %v\n", reflect.DeepEqual(input, output))
	fmt.Printf("error: %v\n", err)
	// Output:
	// equal: true
	// error: <nil>
}
