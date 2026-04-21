package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func printHelp() {
	fmt.Println("Usage: wake-on-lan -m <MAC_ADDRESS> [-t TARGET_IP]")
	fmt.Println()
	fmt.Println("Send a Wake-on-LAN magic packet to wake up a remote machine.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -m, --mac    MAC address of target machine (required)")
	fmt.Println("              Supported formats: 00:11:22:33:44:55 | 00-11-22-33-44-55 | 001122334455")
	fmt.Println("  -t, --target IP address to send packet to (optional, default: 255.255.255.255 broadcast)")
	fmt.Println("              Use unicast IP when broadcast is blocked on the network.")
	fmt.Println("  -h, --help   Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  wake-on-lan -m 00:11:22:33:44:55")
	fmt.Println("  wake-on-lan --mac 00:11:22:33:44:55 --target 192.168.0.198")
}

func main() {
	// Check for help flag before parsing (avoids "required flag" errors)
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			printHelp()
			os.Exit(0)
		}
	}

	macFlag := flag.String("m", "", "MAC address of target machine")
	macFlagLong := flag.String("mac", "", "MAC address of target machine (long form)")
	targetFlag := flag.String("t", "", "Target IP address")
	targetFlagLong := flag.String("target", "", "Target IP address (long form)")
	flag.Parse()

	if *macFlag == "" && *macFlagLong == "" {
		fmt.Fprintf(os.Stderr, "Error: -m/--mac is required\n")
		os.Exit(1)
	}

	macStr := *macFlag
	if macStr == "" {
		macStr = *macFlagLong
	}

	targetIP := "255.255.255.255"
	if *targetFlag != "" {
		targetIP = *targetFlag
	} else if *targetFlagLong != "" {
		targetIP = *targetFlagLong
	}

	mac, err := parseMAC(macStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print which broadcast implementation is active
	fmt.Printf("[%s] broadcast implementation: ", runtime.GOOS)
	if runtime.GOOS == "windows" {
		fmt.Println("net.DialUDP")
	} else {
		fmt.Println("syscall.Socket + SO_BROADCAST")
	}

	usedIP, err := sendWOL(mac, targetIP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if usedIP == "255.255.255.255" {
		fmt.Printf("Wake-on-LAN packet sent via broadcast (%s) for MAC %s\n", usedIP, macStr)
	} else {
		fmt.Printf("Wake-on-LAN packet sent to %s (MAC: %s)\n", usedIP, macStr)
	}
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
// Supports formats: "00:11:22:33:44:55", "00-11-22-33:44-55", "001122334455"
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

// createMagicPacket creates a Wake-on-LAN magic packet for the given MAC address.
// Format: 6 bytes of 0xFF followed by 16 repetitions of the MAC address.
func createMagicPacket(mac []byte) []byte {
	packet := make([]byte, 102)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 0; i < 16; i++ {
		copy(packet[6+i*6:], mac)
	}
	return packet
}

// sendWOL sends a Wake-on-LAN magic packet to the specified MAC address.
// Returns the actual IP address the packet was sent to (may differ from targetIP on fallback).
func sendWOL(mac []byte, targetIP string) (string, error) {
	packet := createMagicPacket(mac)

	err := sendWolBroadcast(packet, targetIP, 9)
	if err == nil {
		return targetIP, nil
	}
	// Fallback to standard broadcast address
	err = sendWolBroadcast(packet, "255.255.255.255", 9)
	return "255.255.255.255", err
}
