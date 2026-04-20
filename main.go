package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("wake-on-lan v1.0.0")
}

// parseHexByte parses a 2-character hex string into a byte.
func parseHexByte(s string) (byte, error) {
	var result byte
	for i := 0; i < 2; i++ {
		c := s[i]
		var val byte
		switch {
		case c >= '0' && c <= '9':
			val = c - '0'
		case c >= 'a' && c <= 'f':
			val = c - 'a' + 10
		case c >= 'A' && c <= 'F':
			val = c - 'A' + 10
		default:
			return 0, fmt.Errorf("invalid hex character: %c", c)
		}
		result = result<<4 | val
	}
	return result, nil
}

// parseMAC parses a MAC address string and returns its byte representation.
// Supports formats: "00:11:22:33:44:55", "00-11-22-33-44-55", "001122334455"
func parseMAC(macStr string) ([]byte, error) {
	cleaned := strings.ReplaceAll(macStr, ":", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	if len(cleaned) != 12 {
		return nil, fmt.Errorf("invalid MAC address format: %s", macStr)
	}

	mac := make([]byte, 6)
	for i := 0; i < 6; i++ {
		byteStr := cleaned[i*2 : i*2+2]
		b, err := parseHexByte(byteStr)
		if err != nil {
			return nil, fmt.Errorf("invalid MAC address: %s", macStr)
		}
		mac[i] = b
	}

	return mac, nil
}
