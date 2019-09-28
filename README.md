# compints

[![GoDoc](https://godoc.org/github.com/shibukawa/compints?status.svg)](https://godoc.org/github.com/shibukawa/compints)

compints is a package that provides compression/decompression algorithm for integer list.

## Functions (Simple)

* ``Compress(input []uint32, diff bool)  (output, cbytes []byte, elementCount int)``
* ``Decompress(input, cbytes []byte, elementCount int, diff bool) ([]uint32, error)``

``Compress`` returns three result. All result are important when calling ``Decompress``.
It is not useful than ``CompressToBytes``, but you can store these values in block devices when update the result.

diff flag makes result more efficient if input integer slice are sorted in ascend order.

```go
var input = []uint32{1, 2, 3, 4, 5, 6, 6, 8, 10}

compressed, cbytes, length := compints.Compress(input, true)
output, err := compints.Decompress(compressed, cbytes, length, true)
// output is as same as input
```

## Functions (Easy)

* ``CompressToBytes(input []uint32, diff bool) []byte``
* ``DecompressFromBytes(input, []byte, diff bool) ([]uint32, error)``

``CompressToBytes`` returns one result. Just pass this result to ``DecompressFromBytes``, you can get decompressed data.

diff flag makes result more efficient if input integer slice are sorted in ascend order.

```go
var input = []uint32{1, 2, 3, 4, 5, 6, 6, 8, 10}

compressed := compints.CompressToBytes(input, true)
output, err := compints.DecompressFromBytes(compressed, true)
// output is as same as input
```

## License

Apache 2