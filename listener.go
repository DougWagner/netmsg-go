package main

import (
	"fmt"
	"net"
)

func ListenOnUDP(addr net.IP, port int, name string) {
	udp, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		fmt.Println(err)
		return
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

func Uint32(b []byte) (n uint32) {
	n += uint32(b[3] << 24)
	n += uint32(b[2] << 16)
	n += uint32(b[1] << 8)
	n += uint32(b[0])
	return
}

func ListenOnTCP(addr net.IP, port int) {
	tcp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		fmt.Println(err)
		return
	}
	tcplisten, err := net.ListenTCP("tcp", tcp)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := tcplisten.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			return
		}
		sizebuf := make([]byte, 4)
		_, err = conn.Read(sizebuf)
		if err != nil {
			fmt.Println(err)
			return
		}
		size := Uint32(sizebuf)
		msgbuf := make([]byte, size)
		_, err = conn.Read(msgbuf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(msgbuf))
		conn.Close()
	}
}
