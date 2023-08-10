package osservice

import (
	"context"
	"fmt"
	"time"

	"osservice/osproto"
	"osservice/pkg/wifi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	_ Wifi = (*WifiClient)(nil)
)

type WifiClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  osproto.WifiClient
}

func NewWifiClient(addr string, port int, timeout time.Duration) (*WifiClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &WifiClient{
		timeout: timeout,
		conn:    conn,
		client:  osproto.NewWifiClient(conn),
	}, nil
}

func (w *WifiClient) ctx() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(w.timeout))
}
func (w *WifiClient) APs() ([]wifi.AP, error) {
	ctx, cancel := w.ctx()
	defer cancel()

	wAps, err := w.client.APs(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	aps := make([]wifi.AP, 0, len(wAps.Ap))
	for _, ap := range wAps.Ap {
		aps = append(aps, wifi.AP{
			ID:   int(ap.GetId()),
			SSID: ap.GetSsid().GetValue(),
		})
	}

	return aps, nil
}

func (w *WifiClient) Connected() (wifi.Status, error) {
	ctx, cancel := w.ctx()
	defer cancel()

	conn, err := w.client.Connected(ctx, &emptypb.Empty{})
	if err != nil {
		return wifi.Status{}, err
	}
	return wifi.Status{
		Connected: conn.Connected.GetValue(),
		SSID:      conn.Ssid.GetValue(),
	}, nil
}

func (w *WifiClient) Disconnect() error {
	ctx, cancel := w.ctx()
	defer cancel()

	_, err := w.client.Disconnect(ctx, &emptypb.Empty{})
	return err
}

func (w *WifiClient) Connect(n wifi.Network) error {
	wNet := osproto.Network{
		Ap: &osproto.AP{
			Id:   int32(n.ID),
			Ssid: wrapperspb.String(n.SSID),
		},
		Password: wrapperspb.String(n.Password),
	}

	ctx, cancel := w.ctx()
	defer cancel()

	_, err := w.client.Connect(ctx, &wNet)
	return err
}
