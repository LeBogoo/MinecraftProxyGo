package utils

import (
	"bufio"
)

func ReadUShort(reader *bufio.Reader) (int, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	secondByte, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return int(firstByte)<<8 | int(secondByte), nil
}

func ToUShort(value int) []byte {
	return []byte{byte(value >> 8), byte(value)}
}
