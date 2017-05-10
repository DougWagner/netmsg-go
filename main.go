package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type GetIPError struct {
	msg string
}

func (e *GetIPError) Error() string {
	return fmt.Sprint(e.msg)
}

func getClientIPNet() (ip *net.IPNet, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ip, err
	}
	for _, i := range ifaces {
		if strings.Contains(i.Flags.String(), "broadcast") {
			addrs, err := i.Addrs()
			if err != nil {
				return ip, err
			}
			for _, j := range addrs {
				ip, ok := j.(*net.IPNet)
				if !ok {
					return ip, &GetIPError{"addr in interface address list is not type *IPNet please report to developer"}
				}
				if ip.IP.To4()!= nil {
					return ip, nil
				}
			}
		}
	}
	return ip, &GetIPError{"Could not obtain IP Address"}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: netmsg-go <username>")
		return
	}
	const port = 34512
	ip, err := getClientIPNet()
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(ip.String())
	go ListenOnUDP(ip.IP.To4(), port, os.Args[1])
	go ListenOnTCP(ip.IP.To4(), port)
	fmt.Println("Listening for messages")
	Broadcast(ip, port, "doug") // for testing - broadcast will belong in main loop
	for {}
}
