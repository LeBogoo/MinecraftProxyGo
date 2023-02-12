package utils

import (
	"bufio"
)

func ReadLong(reader *bufio.Reader) (int64, error) {
	// read an 8 byte long number
	var long int64
	for i := 0; i < 8; i++ {
		byte, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		long |= int64(byte) << uint(8*i)
	}

	return long, nil
}

func ToLong(value int64) []byte {
	// write an 8 byte long number
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(value >> uint(8*i))
	}

	return bytes
}
