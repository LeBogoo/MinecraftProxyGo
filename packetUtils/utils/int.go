package utils

import (
	"bufio"
)

func ReadInt(reader *bufio.Reader) (int, error) {
	packetLength, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return int(packetLength), nil
}

func ReadUInt(reader *bufio.Reader) (uint, error) {
	packetLength, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return uint(packetLength), nil
}

func ToInt(value int) []byte {
	return []byte{byte(value)}
}

func ToUInt(value uint) []byte {
	return []byte{byte(value)}
}
