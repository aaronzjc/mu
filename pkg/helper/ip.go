package helper

import (
	"net"
	"os"
)

func LocalHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func LocalAddr() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
			case *net.IPAddr:
				ip = v.IP
			}
		}
	}
	return ip.String()
}
