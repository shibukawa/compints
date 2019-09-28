package compints

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
	"testing/quick"
)

func TestCompression(t *testing.T) {
	input := []uint32{1, 2, 3, 4, 5, 6, 6, 8, 10}
	compressed, cbytes, length := Compress(input, true)
	t.Log(compressed, len(compressed))
	output, err := Decompress(compressed, cbytes, length, true)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, input, output)
}

func TestRandom(t *testing.T) {
	cycle := func(input []uint32) bool {
		compressed, cbytes, length := Compress(input, false)
		output, err := Decompress(compressed, cbytes, length, false)
		if err != nil {
			return false
		}
		//t.Logf("inputSize = %d (%d bytes) outputSize = %d + %d bytes", len(input), len(input) * 4, len(compressed), len(cbytes))
		if len(input) != len(output) {
			return false
		}
		return reflect.DeepEqual(input, output)
	}
	if err := quick.Check(cycle, nil); err != nil {
		t.Fatal(err)
	}
}

func TestRandomSorted(t *testing.T) {
	cycle := func(input []uint32) bool {
		sort.Slice(input, func(i, j int) bool {
			return input[i] < input[j]
		})
		compressed, cbytes, length := Compress(input, true)
		output, err := Decompress(compressed, cbytes, length, true)
		if err != nil {
			return false
		}
		//t.Logf("inputSize = %d (%d bytes) outputSize = %d + %d bytes", len(input), len(input) * 4, len(compressed), len(cbytes))
		if len(input) != len(output) {
			return false
		}
		return reflect.DeepEqual(input, output)
	}
	if err := quick.Check(cycle, nil); err != nil {
		t.Fatal(err)
	}
}

func TestRandomToBytes(t *testing.T) {
	cycle := func(input []uint32) bool {
		compressed := CompressToBytes(input, false)
		output, err := DecompressFromBytes(compressed, false)
		if err != nil {
			return false
		}
		//t.Logf("inputSize = %d (%d bytes) outputSize = %d + %d bytes", len(input), len(input) * 4, len(compressed), len(cbytes))
		if len(input) != len(output) {
			return false
		}
		return reflect.DeepEqual(input, output)
	}
	if err := quick.Check(cycle, nil); err != nil {
		t.Fatal(err)
	}
}

func TestRandomSortedToBytes(t *testing.T) {
	cycle := func(input []uint32) bool {
		sort.Slice(input, func(i, j int) bool {
			return input[i] < input[j]
		})
		compressed := CompressToBytes(input, false)
		output, err := DecompressFromBytes(compressed, false)
		if err != nil {
			return false
		}
		//t.Logf("inputSize = %d (%d bytes) outputSize = %d + %d bytes", len(input), len(input) * 4, len(compressed), len(cbytes))
		if len(input) != len(output) {
			return false
		}
		return reflect.DeepEqual(input, output)
	}
	if err := quick.Check(cycle, nil); err != nil {
		t.Fatal(err)
	}
}
