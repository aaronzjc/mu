package helper

import (
	"net"
	"net/http"
	"os"
	"strings"
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

func ClientIp(r *http.Request) string {
	ip := r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if res := net.ParseIP(i); res != nil {
			return res.String()
		}
	}
	ip = r.Header.Get("X-Real-IP")
	if res := net.ParseIP(ip); res != nil {
		return res.String()
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
