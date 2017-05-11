package main

import (
	"fmt"
	"net"
)

func uint32ToB(n uint32) (b []byte) {
	b = make([]byte, 4)
	b[0] = byte(n & 255)
	b[1] = byte((n >> 8) & 255)
	b[2] = byte((n >> 16) & 255)
	b[3] = byte((n >> 24) & 255)
	return
}

func sendMessage(localIP, targetIP net.IP, port int, msg string) error {
	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:", localIP))
	if err != nil {
		return err
	}
	targetAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", targetIP, port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", localAddr, targetAddr)
	if err != nil {
		return err
	}
	msgBytes := []byte(msg)
	sizeBytes := uint32ToB(uint32(len(msgBytes)))
	var buffer []byte
	buffer = append(buffer, sizeBytes...)
	buffer = append(buffer, msgBytes...)
	_, err = conn.Write(buffer)
	if err != nil {
		conn.Close()
		return err
	}
	conn.Close()
	return nil
}
