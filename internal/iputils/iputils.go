package iputils

import (
	"errors"
	"net"
)

func LocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func CheckIP(ip string, cidr string) (bool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}
	ipaddr := net.ParseIP(ip)
	if ipaddr == nil {
		return false, errors.New("can't parse ip")
	}
	return ipNet.Contains(ipaddr), nil
}
