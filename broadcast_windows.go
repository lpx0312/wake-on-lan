//go:build windows

package main

import (
	"net"
)

// enableBroadcast enables SO_BROADCAST on the UDP connection
func enableBroadcast(conn *net.UDPConn) error {
	return nil
}

// sendWolBroadcast sends WOL packet using net.DialUDP on Windows
func sendWolBroadcast(packet []byte, targetIP string, targetPort int) error {
	// Parse target IP
	ip := net.ParseIP(targetIP)
	if ip == nil {
		return &net.ParseError{Type: "IP address", Text: targetIP}
	}

	// Create UDP address
	addr := &net.UDPAddr{IP: ip, Port: targetPort}

	// Dial UDP - this sets up SO_BROADCAST automatically on Windows
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send packet
	_, err = conn.Write(packet)
	return err
}