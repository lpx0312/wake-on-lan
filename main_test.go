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

func TestCreateMagicPacket(t *testing.T) {
	mac := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	packet := createMagicPacket(mac)

	// 检查长度
	if len(packet) != 102 {
		t.Errorf("packet length = %d, expected 102", len(packet))
	}

	// 检查前6字节是 0xFF
	for i := 0; i < 6; i++ {
		if packet[i] != 0xFF {
			t.Errorf("packet[%d] = %x, expected 0xFF", i, packet[i])
		}
	}

	// 检查后96字节是16次MAC重复
	for i := 0; i < 16; i++ {
		for j := 0; j < 6; j++ {
			if packet[6+i*6+j] != mac[j] {
				t.Errorf("packet[%d] = %x, expected %x", 6+i*6+j, packet[6+i*6+j], mac[j])
			}
		}
	}
}
