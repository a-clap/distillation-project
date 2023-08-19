package osservice

import (
	"context"
	"fmt"
	"time"

	"osservice/osproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ Net = (*NetClient)(nil)

type NetClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  osproto.NetClient
}

func NewNetClient(addr string, port int, timeout time.Duration) (*NetClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &NetClient{
		timeout: timeout,
		conn:    conn,
		client:  osproto.NewNetClient(conn),
	}, nil
}

func (n *NetClient) ListInterfaces() []NetInterface {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(n.timeout))
	defer cancel()

	nets, err := n.client.ListInterfaces(ctx, &emptypb.Empty{})
	if err != nil {
		return nil
	}

	osNets := make([]NetInterface, 0, len(nets.Interfaces))
	for _, net := range nets.Interfaces {
		osNets = append(osNets, NetInterface{
			Name:     net.Name.GetValue(),
			IPAddrV4: net.Ipaddr.GetValue(),
		})
	}

	return osNets
}
