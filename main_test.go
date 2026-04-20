package main

import (
	"testing"
)

func TestParseMAC(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
		wantErr  bool
	}{
		{"00:11:22:33:44:55", []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, false},
		{"00-11-22-33-44-55", []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, false},
		{"001122334455", []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, false},
		{"invalid", nil, true},
		{"00:11:22:33:44", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseMAC(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMAC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(result) != string(tt.expected) {
				t.Errorf("parseMAC() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
