package main

import (
	"fmt"
	"net"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CreateHost(_ip, _port string) string {
	testIp := net.ParseIP(_ip)
	if testIp.To4() == nil && useDynamicPort == "false" {
		return _ip
	} else {
		return fmt.Sprintf("%s:%s", _ip, _port)
	}
}
