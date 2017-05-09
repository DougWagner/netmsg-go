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

func getClientIPV4() (ip net.IP, err error) {
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
				addr, ok := j.(*net.IPNet)
				if !ok {
					return ip, &GetIPError{"addr in interface address list is not type *IPNet please report to developer"}
				}
				if ip = addr.IP.To4(); ip != nil {
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
	ip, err := getClientIPV4()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ip.String())
	go ListenOnUDP(ip, 34512, os.Args[1])
	for {}
}
