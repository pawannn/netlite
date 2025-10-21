package pkg

import (
	"net"
)

func GetIP() (string, NetliteErr) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", NetliteErr{
			ClientMessage: "Error getting interface addresses",
			Error:         err,
		}
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), NoErr
			}
		}
	}

	return "", NetliteErr{
		ClientMessage: "Unable to get IP of the machine",
		Error:         nil,
	}
}
