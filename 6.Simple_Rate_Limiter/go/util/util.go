package util

import (
	"net"
	"net/http"
	"strings"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-IP")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
		if IPAddress != "" {
			IPAddress = strings.Split(IPAddress, ",")[0]
		}
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
		if host, _, err := net.SplitHostPort(IPAddress); err == nil {
			IPAddress = host
		}
	}
	return IPAddress
}
