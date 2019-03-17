package utils

import (
	"errors"
	"fmt"
	"net"
)

var (
	externalIP string
)

// nolint:gochecknoinits
func init() {
	var err error

	if externalIP, err = getExternalIP(); err != nil {
		panic(fmt.Errorf("get external ip, reason=%v", err))
	}
}

// GetExternalIP return the first available external ip address
func GetExternalIP() string {
	return externalIP
}

func getExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("not connected to the network")
}

// IsErrClosing return true if err indicates socket is closing
func IsErrClosing(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		err = opErr.Err
	}
	// nolint:stylecheck
	return "use of closed network connection" == err.Error()
}
