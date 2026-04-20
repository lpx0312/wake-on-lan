//go:build windows
// +build windows

package main

import (
	"net"
	"syscall"
)

// enableBroadcast enables SO_BROADCAST on a UDPConn (Windows)
func enableBroadcast(conn *net.UDPConn) error {
	file, err := conn.File()
	if err != nil {
		return err
	}
	defer file.Close()

	err = syscall.SetsockoptInt(syscall.Handle(file.Fd()), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	return err
}
