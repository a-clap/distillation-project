package osservice

import (
	"context"

	"osservice/osproto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type NetInterface struct {
	Name     string
	IPAddrV4 string
}

type Net interface {
	ListInterfaces() []NetInterface
}

func (o *Os) ListInterfaces(context.Context, *empty.Empty) (*osproto.Interfaces, error) {
	nets := o.net.ListInterfaces()
	osNets := make([]*osproto.Interface, 0, len(nets))
	for _, net := range nets {
		osNets = append(osNets, &osproto.Interface{
			Name:   wrapperspb.String(net.Name),
			Ipaddr: wrapperspb.String(net.IPAddrV4),
		})
	}
	return &osproto.Interfaces{Interfaces: osNets}, nil
}
