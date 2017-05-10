package main

import (
	"fmt"
	"net"
	"time"
)


func Broadcast(ip *net.IPNet, port int, uname string) {
	broadcastBytes := make([]byte, 4)
	for i, m := range ip.Mask {
		broadcastBytes[i] = ip.IP.To4()[i] | (m ^ 255)
	}
	broadcastIP := net.IP(broadcastBytes)
	bcastUDPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", broadcastIP.To4(), port))
	if err != nil {
		fmt.Println(err)
		return
	}
	lUDPAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:", ip.IP.To4()))
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenUDP("udp", lUDPAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	unameBytes := []byte(uname)
	response := make([]byte, 4)
	conn.WriteToUDP(unameBytes, bcastUDPAddr)
	conn.SetReadDeadline(time.Now().Add(time.Duration(1) * time.Second))
	_, _, err = conn.ReadFromUDP(response)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}
	conn.Close()
	fmt.Println(response)
}
