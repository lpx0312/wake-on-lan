//go:build !windows
// +build !windows

package main

import "net"

// enableBroadcast enables SO_BROADCAST on a UDPConn (Unix - no-op, Unix allows broadcast by default)
func enableBroadcast(conn *net.UDPConn) error {
	return nil
}
