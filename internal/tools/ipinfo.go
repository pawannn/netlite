package tools

import (
	"net"

	"github.com/pawannn/netlite/pkg"
)

type InterfaceInfo struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	MAC         string   `json:"mac"`
	IPv4        []string `json:"ipv4"`
	IPv6        []string `json:"ipv6"`
	IsUp        bool     `json:"is_up"`
	IsLoopback  bool     `json:"is_loopback"`
	IsMulticast bool     `json:"is_multicast"`
}

func IfInfo() ([]InterfaceInfo, pkg.NetliteErr) {
	var infos []InterfaceInfo
	ifaces, err := net.Interfaces()
	if err != nil {
		return infos, pkg.NetliteErr{ClientMessage: "failed to list interfaces", Error: err}
	}

	for _, iface := range ifaces {
		info := InterfaceInfo{
			Name:        iface.Name,
			MAC:         iface.HardwareAddr.String(),
			IsUp:        iface.Flags&net.FlagUp != 0,
			IsLoopback:  iface.Flags&net.FlagLoopback != 0,
			IsMulticast: iface.Flags&net.FlagMulticast != 0,
		}

		if info.IsLoopback {
			info.Type = "Loopback"
		} else if iface.Flags&net.FlagBroadcast != 0 {
			info.Type = "Wi-Fi / Ethernet"
		} else {
			info.Type = "Unknown"
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, _ := net.ParseCIDR(addr.String())
			if ip == nil {
				continue
			}
			if ip.To4() != nil {
				info.IPv4 = append(info.IPv4, ip.String())
			} else {
				info.IPv6 = append(info.IPv6, ip.String())
			}
		}

		infos = append(infos, info)
	}

	return infos, pkg.NoErr
}
