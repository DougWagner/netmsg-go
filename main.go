// netmsg-go is a simple messaging application for local network
// using TCP and UDP network protocols.
// Usage: netmsg-go <username>
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type getIPError struct {
	msg string
}

func (e *getIPError) Error() string {
	return fmt.Sprint(e.msg)
}

func getClientIPNet() (ip *net.IPNet, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ip, err
	}
	for _, i := range ifaces {
		if strings.Contains(i.Flags.String(), "broadcast") && strings.Contains(i.Flags.String(), "up") {
			addrs, err := i.Addrs()
			if err != nil {
				return ip, err
			}
			for _, j := range addrs {
				ip, ok := j.(*net.IPNet)
				if !ok {
					return ip, &getIPError{"addr in interface address list is not type *IPNet please report to developer"}
				}
				if ip.IP.To4() != nil {
					return ip, nil
				}
			}
		}
	}
	return ip, &getIPError{"Could not obtain IP Address"}
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
	fmt.Println("send command syntax: \"send username message contents\"")
	fmt.Println("you do not need to put quotes around your message")
	fmt.Println("type \"exit\" to quit")
	go listenOnUDP(ip.IP.To4(), port, os.Args[1])
	go listenOnTCP(ip.IP.To4(), port)
	fmt.Println("Listening for messages")
	for {
		in := bufio.NewReader(os.Stdin)
		input, err := in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		input = strings.TrimSuffix(input, "\n")
		splitinput := strings.Split(input, " ")
		if strings.ToLower(splitinput[0]) == "send" {
			if len(splitinput) < 3 {
				fmt.Println("Invalid send command")
				continue
			}
			message := strings.Trim(strings.Join(splitinput[2:], " "), " ")
			targetIP, err := broadcast(ip, port, splitinput[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = sendMessage(ip.IP.To4(), *targetIP, port, fmt.Sprintf("%v: %v", os.Args[1], message))
			if err != nil {
				fmt.Println(err)
			}
		} else if strings.ToLower(splitinput[0]) == "exit" {
			return
		} else {
			fmt.Println("Invalid Command")
		}
	}
}
