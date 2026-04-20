//go:build !windows

package main

import (
	"net"
	"syscall"
)

// enableBroadcast enables SO_BROADCAST - no-op on Unix
func enableBroadcast(conn *net.UDPConn) error {
	return nil
}

// sendWolBroadcast sends WOL packet using raw syscall on Unix
func sendWolBroadcast(packet []byte, targetIP string, targetPort int) error {
	// Create UDP socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	// Enable broadcast
	const SOL_SOCKET = 0xffff
	const SO_BROADCAST = 0x20
	err = syscall.SetsockoptInt(fd, SOL_SOCKET, SO_BROADCAST, 1)
	if err != nil {
		return err
	}

	// Parse target IP
	ip := net.ParseIP(targetIP)
	if ip == nil {
		return &net.ParseError{Type: "IP address", Text: targetIP}
	}

	// Create sockaddr
	addr := syscall.SockaddrInet4{Port: targetPort}
	copy(addr.Addr[:], ip.To4())

	// Send
	err = syscall.Sendto(fd, packet, 0, &addr)
	if err != nil {
		return err
	}
	return nil
}
