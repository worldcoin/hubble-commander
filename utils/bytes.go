package utils

// Pad bytes slice with 0s from left to size length
// Returns original bytes slice in case len(slice) >= size
func PadLeft(bytes []byte, size int) []byte {
	l := len(bytes)
	if l >= size {
		return bytes
	}
	paddedBytes := make([]byte, size)
	copy(paddedBytes[size-l:], bytes)
	return paddedBytes
}

func ByteSliceTo32ByteArray(source []byte) [32]byte {
	var target [32]byte
	copy(target[:], source)

	return target
}
