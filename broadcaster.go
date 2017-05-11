package main

import (
	"fmt"
	"net"
	"time"
)

type noResponseError struct {
	msg string
}

func (e *noResponseError) Error() string {
	return fmt.Sprint(e.msg)
}

func broadcast(ip *net.IPNet, port int, uname string) (*net.IP, error) {
	broadcastBytes := make([]byte, 4)
	for i, m := range ip.Mask {
		broadcastBytes[i] = ip.IP.To4()[i] | (m ^ 255)
	}
	broadcastIP := net.IP(broadcastBytes)
	bcastUDPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", broadcastIP.To4(), port))
	if err != nil {
		return nil, err
	}
	lUDPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:", ip.IP.To4()))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", lUDPAddr)
	if err != nil {
		return nil, err
	}
	unameBytes := []byte(uname)
	response := make([]byte, 4)
	conn.WriteToUDP(unameBytes, bcastUDPAddr)
	conn.SetReadDeadline(time.Now().Add(time.Duration(1) * time.Second))
	_, _, err = conn.ReadFromUDP(response)
	if err != nil {
		conn.Close()
		return nil, &noResponseError{"No response received"}
	}
	conn.Close()
	returnIP := net.IP(response)
	return &returnIP, nil
}
