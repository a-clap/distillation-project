package osservice

import (
	"errors"
	"net"
)

var ErrNoIPV4Address = errors.New("interface doesn't have IPv4 address")

var (
	_ Net = netOs{}
)

type netOs struct {
}

func getInterfaceIPv4Addr(interfaceName string) (string, error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
		err      error
	)

	if ief, err = net.InterfaceByName(interfaceName); err != nil {
		return "", err
	}

	if addrs, err = ief.Addrs(); err != nil {
		return "", err
	}

	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			return ipv4Addr.String(), nil
		}
	}

	return "", ErrNoIPV4Address
}

func (netOs) ListInterfaces() []NetInterface {
	var netInterfaces []NetInterface

	for _, name := range []string{"eth0", "wlan0"} {
		if ipv4, err := getInterfaceIPv4Addr(name); err == nil {
			netInterfaces = append(netInterfaces, NetInterface{
				Name:     name,
				IPAddrV4: ipv4,
			})
		}
	}
	return netInterfaces
}
