package utils

import (
	"bufio"
)

func ReadBool(reader *bufio.Reader) (bool, error) {
	packetLength, err := reader.ReadByte()
	if err != nil {
		return false, err
	}

	return packetLength == 0x01, nil
}

func ToBool(value bool) []byte {
	if value {
		return []byte{0x01}
	}

	return []byte{0x00}
}
