package utils

import (
	"os"
	"strconv"
)

func GetUint16OrDefault(name string, default_value uint16) uint16 {
	value := os.Getenv(name)

	if value == "" {
		return default_value
	}

	parsed, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return default_value
	}

	return uint16(parsed)
}

func GetStringOrDefault(name string, default_value string) string {
	value := os.Getenv(name)

	if value == "" {
		return default_value
	}

	return value
}
