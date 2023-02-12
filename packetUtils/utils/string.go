package utils

import (
	"bufio"
)

func ReadString(reader *bufio.Reader, length int) (string, error) {
	bytes := make([]byte, length)
	_, err := reader.Read(bytes)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ToString(value string) []byte {
	return []byte(value)
}
