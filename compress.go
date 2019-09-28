// package compints provides useful wrapper for nelz9999's StreamVByte algorithm.
//
// Input is just slice of uint32. Slice size is flexible.
//
// There are two pairs of functions grouped by output data:
//
//   * Compress/Decompress: returns compressed data, control bytes, element count
//   * CompressToBytes/DecompressFromBytes: returns just []byte
//
// If just serialize the result, CompressToBytes/DecompressFromBytes are useful.
//
// If you stores the value into block devices and keeps flexibility to modify, Compress/Decompress
// allow you to implement more efficient code.
//
// Both functions has mode "diff". If diff==true, my code assumes input data is sorted ascendant order.
package compints

import (
	"bytes"
	"encoding/binary"

	"github.com/nelz9999/stream-vbyte-go/svb"
)

func count(element int) (complete, total int) {
	complete = element / 4
	if element%4 != 0 {
		total = complete + 1
	} else {
		total = complete
	}
	return
}

// Compress receives []uint32 slice and compress into compress data, control bytes, element count.
//
// These three output are required by Decompress.
//
// If you needs more simple method, use CompressToBytes.
func Compress(input []uint32, diff bool) (output, cbytes []byte, elementCount int) {
	completeBlockCount, totalBlockCount := count(len(input))
	output = make([]byte, 0, len(input)*4)

	cbytes = make([]byte, 0, totalBlockCount)
	elementCount = len(input)

	tmpBuffer := make([]byte, 16)
	for i := 0; i < completeBlockCount; i++ {
		cb, n := svb.PutU32Block(tmpBuffer, input[i*4:i*4+4], diff)
		output = append(output, tmpBuffer[0:n]...)
		cbytes = append(cbytes, cb)
	}
	remained := len(input) % 4
	if remained != 0 {
		tmpInput := make([]uint32, 4)
		copy(tmpInput, input[4*completeBlockCount:])
		if diff {
			switch remained {
			case 1:
				tmpInput[1] = tmpInput[0]
				fallthrough
			case 2:
				tmpInput[2] = tmpInput[1]
				fallthrough
			case 3:
				tmpInput[3] = tmpInput[2]
			}
		}
		cb, n := svb.PutU32Block(tmpBuffer, tmpInput, diff)
		output = append(output, tmpBuffer[0:n]...)
		cbytes = append(cbytes, cb)
	}
	return
}

// Decompress consumes the result of Compress and returns original []uint32 slice.
func Decompress(input, cbytes []byte, elementCount int, diff bool) (output []uint32, err error) {
	reader := bytes.NewReader(input)
	output = make([]uint32, 0, len(cbytes)*4)
	for _, c := range cbytes {
		nums, err := svb.ReadUint32s(c, reader)
		if err != nil {
			return nil, err
		}
		if diff {
			nums[1] += nums[0]
			nums[2] += nums[1]
			nums[3] += nums[2]
		}
		output = append(output, nums[0], nums[1], nums[2], nums[3])
	}
	return output[0:elementCount], nil
}

// CompressToBytes receives []uint32 slice and compress into []byte.
func CompressToBytes(input []uint32, diff bool) []byte {
	output, cbytes, elementCount := Compress(input, diff)
	buffer := &bytes.Buffer{}
	buffer.Grow(4 + len(cbytes) + len(output))
	binary.Write(buffer, binary.LittleEndian, uint32(elementCount))
	buffer.Write(cbytes)
	buffer.Write(output)
	return buffer.Bytes()
}

// DecompressFromBytes consumes the result of CompressToBytes and returns original []uint32 slice.
func DecompressFromBytes(input []byte, diff bool) ([]uint32, error) {
	reader := bytes.NewReader(input)
	var elementCount uint32
	err := binary.Read(reader, binary.LittleEndian, &elementCount)
	if err != nil {
		return nil, err
	}
	_, totalBlockCount := count(int(elementCount))
	cbytes := input[4 : 4+totalBlockCount]
	compressed := input[4+totalBlockCount:]

	return Decompress(compressed, cbytes, int(elementCount), diff)
}
