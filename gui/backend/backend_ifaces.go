package backend

type NetInterface struct {
	Name   string `json:"name"`
	IPAddr string `json:"ip_addr"`
}

func (b *Backend) ListInterfaces() []NetInterface {
	if b.net == nil {
		return nil
	}
	ifaces := b.net.ListInterfaces()

	netInterfaces := make([]NetInterface, len(ifaces))
	for i, iface := range ifaces {
		netInterfaces[i] = NetInterface{
			Name:   iface.Name,
			IPAddr: iface.IPAddrV4,
		}
	}

	return netInterfaces
}
