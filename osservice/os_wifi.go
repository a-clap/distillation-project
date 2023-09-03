package osservice

import (
	"context"

	"osservice/osproto"
	"osservice/pkg/wifi"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Wifi interface {
	APs() ([]wifi.AP, error)
	Connected() (wifi.Status, error)
	Disconnect() error
	Connect(n wifi.Network) error
}

func (o *Os) APs(context.Context, *empty.Empty) (*osproto.APReplies, error) {
	aps, err := o.wifi.APs()
	if err != nil {
		return nil, err
	}

	pAps := &osproto.APReplies{Ap: make([]*osproto.AP, 0, len(aps))}
	for _, ap := range aps {
		pAps.Ap = append(pAps.Ap, &osproto.AP{
			Id:   int32(ap.ID),
			Ssid: wrapperspb.String(ap.SSID),
		})
	}

	return pAps, nil
}

func (o *Os) Disconnect(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, o.wifi.Disconnect()
}

func (o *Os) Connected(context.Context, *empty.Empty) (*osproto.ConnectedReply, error) {
	conn, err := o.wifi.Connected()
	if err != nil {
		return nil, err
	}
	return &osproto.ConnectedReply{
		Connected: wrapperspb.Bool(conn.Connected),
		Ssid:      wrapperspb.String(conn.SSID),
	}, nil
}

func (o *Os) Connect(_ context.Context, network *osproto.Network) (*empty.Empty, error) {
	net := wifi.Network{
		AP: wifi.AP{
			ID:   int(network.Ap.Id),
			SSID: network.Ap.Ssid.GetValue(),
		},
		Password: network.Password.GetValue(),
	}
	err := o.wifi.Connect(net)

	return &empty.Empty{}, err
}
