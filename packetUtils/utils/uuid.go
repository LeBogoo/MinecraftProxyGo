package utils

import (
	"bufio"
	"encoding/hex"
)

func ReadUUID(reader *bufio.Reader) (string, error) {
	bytes := make([]byte, 16)
	_, err := reader.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
