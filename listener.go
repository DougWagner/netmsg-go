package main

import (
	"fmt"
	"net"
)

func ListenOnUDP(addr net.IP, port int, name string) {
	udp, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.ListenUDP("udp", udp)
	for {
		if err != nil {
			fmt.Println(err)
			return
		}
		buffer := make([]byte, 1024)
		read, returnIP, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		uname := string(buffer[:read])
		fmt.Println(uname, returnIP)
		if uname == name {
			_, err := conn.WriteToUDP(addr, returnIP)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
