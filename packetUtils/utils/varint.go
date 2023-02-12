package utils

import (
	"bufio"
	"fmt"
)

const (
	SEGMENT_BITS = 0x7F
	CONTINUE_BIT = 0x80
)

func ReadVarInt(reader *bufio.Reader) (int, error) {
	value := 0
	position := 0
	var currentByte byte

	for {
		var err error
		currentByte, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		value |= int(currentByte&SEGMENT_BITS) << position

		if (currentByte & CONTINUE_BIT) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, fmt.Errorf("VarInt is too big")
		}
	}

	return value, nil
}

func ToVarInt(value int) []byte {
	var result []byte
	for {
		if (value & ^SEGMENT_BITS) == 0 {
			result = append(result, byte(value))
			break
		}

		result = append(result, byte((value&SEGMENT_BITS)|CONTINUE_BIT))
		value = int(uint32(value) >> 7)
	}

	return result
}
